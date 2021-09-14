package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/ethereum2"
	"github.com/kotalco/api/shared"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ValidatorHandler is Ethereum 2.0 validator client handler
type ValidatorHandler struct{}

// NewValidatorHandler creates a new Ethereum 2.0 validator client handler
func NewValidatorHandler() handlers.Handler {
	return &ValidatorHandler{}
}

// Get gets a single Ethereum 2.0 validator client by name
func (v *ValidatorHandler) Get(c *fiber.Ctx) error {
	validator := c.Locals("validator").(*ethereum2v1alpha1.Validator)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validator": models.FromEthereum2Validator(validator),
	})
}

// List returns all Ethereum 2.0 validator clients
func (v *ValidatorHandler) List(c *fiber.Ctx) error {
	validators := &ethereum2v1alpha1.ValidatorList{}
	validatorModels := []models.Validator{}

	if err := k8s.Client().List(c.Context(), validators); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all validators",
		})
	}

	page := c.Query("page")
	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	start, end := shared.Page(uint(len(validators.Items)), uint(p))
	sort.Slice(validators.Items[:], func(i, j int) bool {
		return validators.Items[i].CreationTimestamp.Before(&validators.Items[j].CreationTimestamp)
	})

	for _, validator := range validators.Items[start:end] {
		validatorModels = append(validatorModels, *models.FromEthereum2Validator(&validator))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validators": validatorModels,
	})
}

// Create creates Ethereum 2.0 validator client from spec
func (v *ValidatorHandler) Create(c *fiber.Ctx) error {
	model := new(models.Validator)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	keystores := []ethereum2v1alpha1.Keystore{}
	for _, keystore := range model.Keystores {
		keystores = append(keystores, ethereum2v1alpha1.Keystore{
			SecretName: keystore.SecretName,
		})
	}

	validator := &ethereum2v1alpha1.Validator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: "default",
		},
		Spec: ethereum2v1alpha1.ValidatorSpec{
			Network:   model.Network,
			Client:    ethereum2v1alpha1.Ethereum2Client(model.Client),
			Keystores: keystores,
		},
	}

	if model.Client == string(ethereum2v1alpha1.PrysmClient) && model.WalletPasswordSecretName != "" {
		validator.Spec.WalletPasswordSecret = model.WalletPasswordSecretName
	}

	if len(model.BeaconEndpoints) != 0 {
		validator.Spec.BeaconEndpoints = model.BeaconEndpoints
	} else {
		validator.Spec.BeaconEndpoints = []string{}
	}

	if os.Getenv("MOCK") == "true" {
		validator.Default()
	}

	if err := k8s.Client().Create(c.Context(), validator); err != nil {
		log.Println(err)
		if errors.IsAlreadyExists(err) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": fmt.Sprintf("validator by name %s already exist", model.Name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create validator",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validator": models.FromEthereum2Validator(validator),
	})
}

// Delete deletes Ethereum 2.0 validator client by name
func (v *ValidatorHandler) Delete(c *fiber.Ctx) error {
	validator := c.Locals("validator").(*ethereum2v1alpha1.Validator)

	if err := k8s.Client().Delete(c.Context(), validator); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete validator by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates Ethereum 2.0 validator client by name from spec
func (v *ValidatorHandler) Update(c *fiber.Ctx) error {
	model := new(models.Validator)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	name := c.Params("name")
	validator := c.Locals("validator").(*ethereum2v1alpha1.Validator)

	if model.WalletPasswordSecretName != "" {
		validator.Spec.WalletPasswordSecret = model.WalletPasswordSecretName
	}

	if len(model.Keystores) != 0 {
		keystores := []ethereum2v1alpha1.Keystore{}
		for _, keystore := range model.Keystores {
			keystores = append(keystores, ethereum2v1alpha1.Keystore{
				SecretName: keystore.SecretName,
			})
		}
		validator.Spec.Keystores = keystores
	}

	if model.Graffiti != "" {
		validator.Spec.Graffiti = model.Graffiti
	}

	if len(model.BeaconEndpoints) != 0 {
		validator.Spec.BeaconEndpoints = model.BeaconEndpoints
	}

	if model.CPU != "" {
		validator.Spec.CPU = model.CPU
	}
	if model.CPULimit != "" {
		validator.Spec.CPULimit = model.CPULimit
	}
	if model.Memory != "" {
		validator.Spec.Memory = model.Memory
	}
	if model.MemoryLimit != "" {
		validator.Spec.MemoryLimit = model.MemoryLimit
	}
	if model.Storage != "" {
		validator.Spec.Storage = model.Storage
	}

	if os.Getenv("MOCK") == "true" {
		validator.Default()
	}

	if err := k8s.Client().Update(c.Context(), validator); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't update validator by name %s", name),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validator": models.FromEthereum2Validator(validator),
	})
}

// Count returns total number of validators
func (pr *ValidatorHandler) Count(c *fiber.Ctx) error {
	validators := &ethereum2v1alpha1.ValidatorList{}
	if err := k8s.Client().List(c.Context(), validators); err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("X-Total-Count", fmt.Sprintf("%d", len(validators.Items)))

	return c.SendStatus(http.StatusOK)
}

// validateValidatorExist validate node by name exist
func validateValidatorExist(c *fiber.Ctx) error {
	name := c.Params("name")
	validator := &ethereum2v1alpha1.Validator{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, validator); err != nil {

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("validator by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get validator by name %s", name),
		})
	}

	c.Locals("validator", validator)

	return c.Next()

}

// Register registers all handlers on the given router
func (v *ValidatorHandler) Register(router fiber.Router) {
	router.Post("/", v.Create)
	router.Head("/", v.Count)
	router.Get("/", v.List)
	router.Get("/:name", validateValidatorExist, v.Get)
	router.Put("/:name", validateValidatorExist, v.Update)
	router.Delete("/:name", validateValidatorExist, v.Delete)
}

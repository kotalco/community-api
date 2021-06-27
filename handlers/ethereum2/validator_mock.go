package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	models "github.com/kotalco/api/models/ethereum2"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ValidatorMockHandler is Ethereum 2.0 mock validator client handler
type ValidatorMockHandler struct{}

var validatorsStore = map[string]*ethereum2v1alpha1.Validator{}

// NewValidatorMockHandler creates a new Ethereum 2.0 mock validator client handler
func NewValidatorMockHandler() handlers.Handler {
	return &ValidatorMockHandler{}
}

// Get gets a single Ethereum 2.0 mock validator client by name
func (p *ValidatorMockHandler) Get(c *fiber.Ctx) error {
	name := c.Params("name")
	validator := validatorsStore[name]
	model := models.FromEthereum2Validator(validator)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validator": model,
	})
}

// List returns all Ethereum 2.0 mock validator clients
func (p *ValidatorMockHandler) List(c *fiber.Ctx) error {
	validators := []models.Validator{}
	for _, validator := range validatorsStore {
		validators = append(validators, *models.FromEthereum2Validator(validator))
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validators": validators,
	})
}

// Create creates Ethereum 2.0 mock validator client from spec
func (p *ValidatorMockHandler) Create(c *fiber.Ctx) error {
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
			Name: model.Name,
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

	validator.Default()

	validatorsStore[model.Name] = validator

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validator": models.FromEthereum2Validator(validator),
	})

}

// Delete deletes Ethereum 2.0 mock validator client by name
func (p *ValidatorMockHandler) Delete(c *fiber.Ctx) error {
	name := c.Params("name")
	delete(validatorsStore, name)
	return c.SendStatus(http.StatusNoContent)
}

// Update updates Ethereum 2.0 mock validator client by name from spec
func (p *ValidatorMockHandler) Update(c *fiber.Ctx) error {
	model := new(models.Validator)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	name := c.Params("name")
	validator := validatorsStore[name]

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

	validator.Default()

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validator": models.FromEthereum2Validator(validator),
	})
}

// validateValidatorMockExist validate validator client by name exist
func validateValidatorMockExist(c *fiber.Ctx) error {
	name := c.Params("name")

	if validatorsStore[name] != nil {
		return c.Next()
	}
	return c.Status(http.StatusNotFound).JSON(map[string]string{
		"error": fmt.Sprintf("validator by name %s doesn't exist", name),
	})
}

// Register registers all handlers on the given router
func (p *ValidatorMockHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", validateValidatorMockExist, p.Get)
	router.Put("/:name", validateValidatorMockExist, p.Update)
	router.Delete("/:name", validateValidatorMockExist, p.Delete)
}

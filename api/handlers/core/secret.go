package handlers

import (
	"fmt"
	"github.com/kotalco/api/api/handlers"
	"github.com/kotalco/api/internal/models/core"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/shared"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretHandler is k8s secret handler
type SecretHandler struct{}

// NewSecretHandler creates a new k8s secret handler
func NewSecretHandler() handlers.Handler {
	return &SecretHandler{}
}

// Get gets a single k8s secret by name
func (s *SecretHandler) Get(c *fiber.Ctx) error {
	secret := c.Locals("secret").(*corev1.Secret)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"secret": models.FromCoreSecret(secret),
	})
}

// List returns all k8s secrets
func (s *SecretHandler) List(c *fiber.Ctx) error {
	secrets := &corev1.SecretList{}
	if err := k8s.Client().List(c.Context(), secrets, client.InNamespace("default"), client.HasLabels{"app.kubernetes.io/created-by"}); err != nil {
		log.Println(err)

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all secrets",
		})
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(secrets.Items)))

	secretModels := []models.Secret{}
	secretType := c.Query("type")
	// default page to 0
	page, _ := strconv.Atoi(c.Query("page"))

	start, end := shared.Page(uint(len(secrets.Items)), uint(page))
	sort.Slice(secrets.Items[:], func(i, j int) bool {
		return secrets.Items[j].CreationTimestamp.Before(&secrets.Items[i].CreationTimestamp)
	})

	for _, secret := range secrets.Items[start:end] {
		keyType := secret.Labels["kotal.io/key-type"]

		if keyType == "" || secretType != "" && keyType != secretType {
			continue
		}
		secretModels = append(secretModels, *models.FromCoreSecret(&secret))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"secrets": secretModels,
	})

}

// Create creates k8s secret from spec
func (s *SecretHandler) Create(c *fiber.Ctx) error {
	model := new(models.Secret)

	if err := c.BodyParser(model); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	t := true
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: "default",
			Labels: map[string]string{
				"kotal.io/key-type":            model.Type,
				"app.kubernetes.io/created-by": "kotal-api",
			},
		},
		StringData: model.Data,
		Immutable:  &t,
	}

	if err := k8s.Client().Create(c.Context(), secret); err != nil {
		log.Println(err)
		if errors.IsAlreadyExists(err) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": fmt.Sprintf("secret by name %s already exist", model.Name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create secret",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"secret": models.FromCoreSecret(secret),
	})
}

// Delete deletes k8s secret by name
func (s *SecretHandler) Delete(c *fiber.Ctx) error {
	secret := c.Locals("secret").(*corev1.Secret)

	if err := k8s.Client().Delete(c.Context(), secret); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete secret by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates k8s secret by name from spec
func (s *SecretHandler) Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Count returns total number of secrets
func (s *SecretHandler) Count(c *fiber.Ctx) error {
	secrets := &corev1.SecretList{}
	if err := k8s.Client().List(c.Context(), secrets, client.InNamespace("default"), client.HasLabels{"kotal.io/key-type"}); err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(secrets.Items)))

	return c.SendStatus(http.StatusOK)
}

// validateSecretExist validate secret by name exist
func validateSecretExist(c *fiber.Ctx) error {
	name := c.Params("name")
	secret := &corev1.Secret{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, secret); err != nil {
		log.Println(err)
		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("secret by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get secret by name %s", name),
		})
	}

	c.Locals("secret", secret)

	return c.Next()

}

// Register registers all handlers on the given router
func (s *SecretHandler) Register(router fiber.Router) {
	router.Post("/", s.Create)
	router.Head("/", s.Count)
	router.Get("/", s.List)
	router.Get("/:name", validateSecretExist, s.Get)
	router.Put("/:name", validateSecretExist, s.Update)
	router.Delete("/:name", validateSecretExist, s.Delete)
}

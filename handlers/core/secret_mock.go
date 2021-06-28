package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	models "github.com/kotalco/api/models/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretMockHandler is k8s secret mock handler
type SecretMockHandler struct{}

// secretsStore is in-memory secrets store
var secretsStore = map[string]*corev1.Secret{}

// NewSecretMockHandler creates a new k8s secret mock handler
func NewSecretMockHandler() handlers.Handler {
	return &SecretMockHandler{}
}

// Get gets a single k8s secret mock by name
func (s *SecretMockHandler) Get(c *fiber.Ctx) error {
	name := c.Params("name")
	secret := secretsStore[name]
	model := models.FromCoreSecret(secret)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"secret": model,
	})
}

// List returns all k8s secret mocks
func (s *SecretMockHandler) List(c *fiber.Ctx) error {
	secrets := []models.Secret{}

	for _, secret := range secretsStore {
		secrets = append(secrets, *models.FromCoreSecret(secret))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"secrets": secrets,
	})
}

// Create creates k8s secret mock from spec
func (s *SecretMockHandler) Create(c *fiber.Ctx) error {
	model := new(models.Secret)

	if err := c.BodyParser(model); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	// check if secret exist with this name
	if secretsStore[model.Name] != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(map[string]string{
			"error": fmt.Sprintf("secret by name %s already exist", model.Name),
		})
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: model.Name,
			Labels: map[string]string{
				"kotal.io/key-type":            model.Type,
				"app.kubernetes.io/created-by": "kotal-api",
			},
		},
		StringData: model.Data,
	}

	secretsStore[model.Name] = secret

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"secret": models.FromCoreSecret(secret),
	})
}

// Delete deletes k8s secret mock by name
func (s *SecretMockHandler) Delete(c *fiber.Ctx) error {
	name := c.Params("name")
	delete(secretsStore, name)
	return c.SendStatus(http.StatusNoContent)
}

// Update updates k8s secret mock by name from spec
func (s *SecretMockHandler) Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// validateSecretExist validate secret by name exist
func validateSecretExist(c *fiber.Ctx) error {
	name := c.Params("name")

	if secretsStore[name] != nil {
		return c.Next()
	}
	return c.Status(http.StatusNotFound).JSON(map[string]string{
		"error": fmt.Sprintf("secret by name %s doesn't exist", c.Params("name")),
	})
}

// Register registers all handlers on the given router
func (s *SecretMockHandler) Register(router fiber.Router) {
	router.Post("/", s.Create)
	router.Get("/", s.List)
	router.Get("/:name", validateSecretExist, s.Get)
	router.Put("/:name", validateSecretExist, s.Update)
	router.Delete("/:name", validateSecretExist, s.Delete)
}

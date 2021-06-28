package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/core"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// SecretHandler is k8s secret handler
type SecretHandler struct{}

// NewSecretHandler creates a new k8s secret handler
func NewSecretHandler() handlers.Handler {
	return &SecretHandler{}
}

// Get gets a single k8s secret by name
func (s *SecretHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a secret")
}

// List returns all k8s secrets
func (s *SecretHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all secrets")
}

// Create creates k8s secret from spec
func (s *SecretHandler) Create(c *fiber.Ctx) error {
	model := new(models.Secret)

	if err := c.BodyParser(model); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

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
	return c.SendString("Delete a secret")
}

// Update updates k8s secret by name from spec
func (s *SecretHandler) Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
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
	router.Get("/", s.List)
	router.Get("/:name", validateSecretExist, s.Get)
	router.Put("/:name", validateSecretExist, s.Update)
	router.Delete("/:name", validateSecretExist, s.Delete)
}

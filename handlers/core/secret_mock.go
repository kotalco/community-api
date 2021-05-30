package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// SecretMockHandler is k8s secret mock handler
type SecretMockHandler struct{}

// NewSecretMockHandler creates a new k8s secret mock handler
func NewSecretMockHandler() handlers.Handler {
	return &SecretMockHandler{}
}

// Get gets a single k8s secret mock by name
func (s *SecretMockHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a mock k8s secret")
}

// List returns all k8s secret mocks
func (s *SecretMockHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all mock k8s secrets")
}

// Create creates k8s secret mock from spec
func (s *SecretMockHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a mock k8s secret")
}

// Delete deletes k8s secret mock by name
func (s *SecretMockHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a mock k8s secret")
}

// Update updates k8s secret mock by name from spec
func (s *SecretMockHandler) Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Register registers all handlers on the given router
func (s *SecretMockHandler) Register(router fiber.Router) {
	router.Post("/", s.Create)
	router.Get("/", s.List)
	router.Get("/:name", s.Get)
	router.Put("/:name", s.Update)
	router.Delete("/:name", s.Delete)
}

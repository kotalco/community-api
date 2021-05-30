package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
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
	return c.SendString("Create a secret")
}

// Delete deletes k8s secret by name
func (s *SecretHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a secret")
}

// Update updates k8s secret by name from spec
func (s *SecretHandler) Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Register registers all handlers on the given router
func (s *SecretHandler) Register(router fiber.Router) {
	router.Post("/", s.Create)
	router.Get("/", s.List)
	router.Get("/:name", s.Get)
	router.Put("/:name", s.Update)
	router.Delete("/:name", s.Delete)
}

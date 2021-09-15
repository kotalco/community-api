package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// StorageClassHandler is k8s storage class handler
type StorageClassHandler struct{}

// NewStorageClassHandler creates a new k8s storage class handler
func NewStorageClassHandler() handlers.Handler {
	return &StorageClassHandler{}
}

// Get gets a single k8s storage class
func (s *StorageClassHandler) Get(c *fiber.Ctx) error {
	return c.SendString("get a single storage class by name")
}

// List returns all k8s storage classes
func (s *StorageClassHandler) List(c *fiber.Ctx) error {
	return c.SendString("list all storage classes")
}

// Create creates k8s storage class from spec
func (s *StorageClassHandler) Create(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Delete deletes k8s storage class by name
func (s *StorageClassHandler) Delete(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Update updates k8s storage class by name from spec
func (s *StorageClassHandler) Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Register registers all handlers on the given router
func (s *StorageClassHandler) Register(router fiber.Router) {
	router.Post("/", s.Create)
	router.Get("/", s.List)
	router.Get("/:name", s.Get)
	router.Put("/:name", s.Update)
	router.Delete("/:name", s.Delete)
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/core"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// StorageClassHandler is k8s storage class handler
type StorageClassHandler struct{}

// NewStorageClassHandler creates a new k8s storage class handler
func NewStorageClassHandler() handlers.Handler {
	return &StorageClassHandler{}
}

// Get gets a single k8s storage class
func (s *StorageClassHandler) Get(c *fiber.Ctx) error {
	name := c.Params("name")
	storageClass := &storagev1.StorageClass{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, storageClass); err != nil {

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("storage class by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get storage class by name %s", name),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"storageClass": models.FromCoreStorageClass(storageClass),
	})
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

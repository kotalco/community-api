package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// BeaconNodeHandler is ethereum 2.0 beacon node handler
type BeaconNodeHandler struct{}

// NewBeaconNodeHandler creates a new ethereum 2.0 beacon node handler
func NewBeaconNodeHandler() handlers.Handler {
	return &BeaconNodeHandler{}
}

// Get gets a single ethereum 2.0 beacon node by name
func (p *BeaconNodeHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a beacon node")
}

// List returns all ethereum 2.0 beacon nodes
func (p *BeaconNodeHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all beacon nodes")
}

// Create creates ethereum 2.0 beacon node from spec
func (p *BeaconNodeHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a beacon node")
}

// Delete deletes ethereum 2.0 beacon node by name
func (p *BeaconNodeHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a beacon node")
}

// Update updates ethereum 2.0 beacon node by name from spec
func (p *BeaconNodeHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a beacon node")
}

// validateBeaconNodeExist validate node by name exist
func validateBeaconNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")
	node := &ethereum2v1alpha1.BeaconNode{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, node); err != nil {

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("beacon node by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get beacon node by name %s", name),
		})
	}

	c.Locals("node", node)

	return c.Next()

}

// Register registers all handlers on the given router
func (p *BeaconNodeHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", validateBeaconNodeExist, p.Get)
	router.Put("/:name", validateBeaconNodeExist, p.Update)
	router.Delete("/:name", validateBeaconNodeExist, p.Delete)
}

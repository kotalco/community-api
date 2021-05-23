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

// BeaconNodeMockHandler is ethereum 2.0 beacon node mock handler
type BeaconNodeMockHandler struct{}

// beaconnodesStore is in-memory beacon nodes store
var beaconnodesStore = map[string]*ethereum2v1alpha1.BeaconNode{}

// NewBeaconNodeMockHandler creates a new ethereum 2.0 beacon node mock handler
func NewBeaconNodeMockHandler() handlers.Handler {
	return &BeaconNodeMockHandler{}
}

// Get gets a single ethereum 2.0 beacon node by name
func (p *BeaconNodeMockHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a mock beacon node")
}

// List returns all ethereum 2.0 beacon nodes
func (p *BeaconNodeMockHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all mock beacon nodes")
}

// Create creates ethereum 2.0 beacon node from spec
func (p *BeaconNodeMockHandler) Create(c *fiber.Ctx) error {
	model := new(models.BeaconNode)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if beaconnodesStore[model.Name] != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("beacon node by name %s already exist", model.Name),
		})
	}

	beaconnodesStore[model.Name] = &ethereum2v1alpha1.BeaconNode{
		ObjectMeta: metav1.ObjectMeta{
			Name: model.Name,
		},
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"beaconnode": model,
	})
}

// Delete deletes ethereum 2.0 beacon node by name
func (p *BeaconNodeMockHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a mock beacon node")
}

// Update updates ethereum 2.0 beacon node by name from spec
func (p *BeaconNodeMockHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a mock beacon node")
}

// Register registers all handlers on the given router
func (p *BeaconNodeMockHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

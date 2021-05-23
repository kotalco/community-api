package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// BeaconNodeMockHandler is ethereum 2.0 beacon node mock handler
type BeaconNodeMockHandler struct{}

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
	return c.SendString("Create a mock beacon node")
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

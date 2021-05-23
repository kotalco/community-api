package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
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

// Register registers all handlers on the given router
func (p *BeaconNodeHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

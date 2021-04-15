package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// PeerMockHandler is IPFS mock peer handler
type PeerMockHandler struct{}

// NewMockPeerHandler creates a new mock IPFS peer handler
func NewPeerMockHandler() handlers.Handler {
	return &PeerMockHandler{}
}

// Get gets a single mock IPFS peer by name
func (p *PeerMockHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a mock peer")
}

// List returns all IPFS mock peers
func (p *PeerMockHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all mock peers")
}

// Create creates IPFS mock peer from spec
func (p *PeerMockHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a mock peer")
}

// Delete deletes IPFS mock peer by name
func (p *PeerMockHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a mock peer")
}

// Update updates IPFS mock peer by name from spec
func (p *PeerMockHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a mock peer")
}

// Register registers all handlers on the given router
func (p *PeerMockHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

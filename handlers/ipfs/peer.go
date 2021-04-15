package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// PeerHandler is IPFS peer handler
type PeerHandler struct{}

// NewPeerHandler creates a new IPFS peer handler
func NewPeerHandler() handlers.Handler {
	return &PeerHandler{}
}

// Get gets a single IPFS peer by name
func (p *PeerHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a peer")
}

// List returns all IPFS peers
func (p *PeerHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all peers")
}

// Create creates IPFS peer from spec
func (p *PeerHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a peer")
}

// Delete deletes IPFS peer by name
func (p *PeerHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a peer")
}

// Update updates IPFS peer by name from spec
func (p *PeerHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a peer")
}

// Register registers all handlers on the given router
func (p *PeerHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

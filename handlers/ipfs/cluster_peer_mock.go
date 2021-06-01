package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// ClusterPeerMockHandler is IPFS peer handler
type ClusterPeerMockHandler struct{}

// NewClusterPeerMockHandler creates a new IPFS cluster peer handler
func NewClusterPeerMockHandler() handlers.Handler {
	return &ClusterPeerMockHandler{}
}

// Get gets a single IPFS cluster peer by name
func (p *ClusterPeerMockHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a mock cluster peer")
}

// List returns all IPFS cluster peers
func (p *ClusterPeerMockHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all mock cluster peers")
}

// Create creates IPFS cluster peer from spec
func (p *ClusterPeerMockHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a mock cluster peer")
}

// Delete deletes IPFS cluster peer by name
func (p *ClusterPeerMockHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a mock cluster peer")
}

// Update updates IPFS cluster peer by name from spec
func (p *ClusterPeerMockHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a mock cluster peer")
}

// Register registers all handlers on the given router
func (p *ClusterPeerMockHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

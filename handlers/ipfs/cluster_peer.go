package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// ClusterPeerHandler is IPFS peer handler
type ClusterPeerHandler struct{}

// NewClusterPeerHandler creates a new IPFS cluster peer handler
func NewClusterPeerHandler() handlers.Handler {
	return &ClusterPeerHandler{}
}

// Get gets a single IPFS cluster peer by name
func (p *ClusterPeerHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a cluster peer")
}

// List returns all IPFS cluster peers
func (p *ClusterPeerHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all cluster peers")
}

// Create creates IPFS cluster peer from spec
func (p *ClusterPeerHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a cluster peer")
}

// Delete deletes IPFS cluster peer by name
func (p *ClusterPeerHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a cluster peer")
}

// Update updates IPFS cluster peer by name from spec
func (p *ClusterPeerHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a cluster peer")
}

// Register registers all handlers on the given router
func (p *ClusterPeerHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

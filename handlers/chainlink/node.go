package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/api/handlers"
	sharedHandlers "github.com/kotalco/api/handlers/shared"
)

// Chainlink node handler
type NodeHandler struct{}

// NewNodeHandler returns new chainlink node handler
func NewNodeHandler() handlers.Handler {
	return &NodeHandler{}
}

// Get returns a single chainlink node by name
func (n *NodeHandler) Get(c *fiber.Ctx) error {
	return nil
}

// Create creates chainlink node from the given spec
func (n *NodeHandler) Create(c *fiber.Ctx) error {
	return nil
}

// Update updates a single chainlink node by name from spec
func (n *NodeHandler) Update(c *fiber.Ctx) error {
	return nil
}

// List returns all chainlink nodes
func (n *NodeHandler) List(c *fiber.Ctx) error {
	return nil
}

// Delete a single chainlink node by name
func (n *NodeHandler) Delete(c *fiber.Ctx) error {
	return nil
}

// Count returns total number of nodes
func (n *NodeHandler) Count(c *fiber.Ctx) error {
	return nil
}

// validateNodeExist validate node by name exist
func validateNodeExist(c *fiber.Ctx) error {
	return c.Next()
}

func (n *NodeHandler) Register(router fiber.Router) {
	router.Post("/", n.Create)
	router.Head("/", n.Count)
	router.Get("/", n.List)
	router.Get("/:name", validateNodeExist, n.Get)
	router.Get("/:name/logs", websocket.New(sharedHandlers.Logger))
	router.Get("/:name/status", websocket.New(sharedHandlers.Status))
	router.Put("/:name", validateNodeExist, n.Update)
	router.Delete("/:name", validateNodeExist, n.Delete)
}

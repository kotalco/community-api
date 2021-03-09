package handlers

import "github.com/gofiber/fiber/v2"

// NodeMockHandler is Ethereum node mock handler
type NodeMockHandler struct{}

// NewNodeMockHandler returns new Ethereum node mock handler
func NewNodeMockHandler() Handler {
	return &NodeMockHandler{}
}

// Get gets a single node
func (e *NodeMockHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get mock node")
}

// List lists all Ethereum nodes
func (e *NodeMockHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all mock nodes")
}

// Create creates a single Ethereum node from spec
func (e *NodeMockHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create mock node")
}

// Delete deletes a single Ethereum node by name
func (e *NodeMockHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete mock node")
}

// Update updates a single node by name from new spec delta
func (e *NodeMockHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update mock node")
}

// Register registers all routes on the given router
func (e *NodeMockHandler) Register(router fiber.Router) {
	router.Post("/", e.Create)
	router.Get("/", e.List)
	router.Get("/:name", e.Get)
	router.Put("/:name", e.Update)
	router.Delete("/:name", e.Delete)
}

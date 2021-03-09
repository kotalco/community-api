package handlers

import "github.com/gofiber/fiber/v2"

// NodeHandler is Ethereum node handler
type NodeHandler struct{}

// NewNodeHandler returns new Ethereum node handler
func NewNodeHandler() Handler {
	return &NodeHandler{}
}

// Get gets a single node
func (e *NodeHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a node")
}

// List lists all Ethereum nodes
func (e *NodeHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all nodes")
}

// Create creates a single Ethereum node from spec
func (e *NodeHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a node")
}

// Delete deletes a single Ethereum node by name
func (e *NodeHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a node")
}

// Update updates a single node by name from new spec delta
func (e *NodeHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a node")
}

// Register registers all routes on the given router
func (e *NodeHandler) Register(router fiber.Router) {
	router.Post("/", e.Create)
	router.Get("/", e.List)
	router.Get("/:name", e.Get)
	router.Put("/:name", e.Update)
	router.Delete("/:name", e.Delete)
}

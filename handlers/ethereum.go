package handlers

import "github.com/gofiber/fiber/v2"

// EthereumHandler is Ethereum handler
type EthereumHandler struct{}

// NewEthereumHandler returns new mock handler
func NewEthereumHandler() Handler {
	return &EthereumHandler{}
}

// Get gets a single node
func (e *EthereumHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a node")
}

// List lists all Ethereum nodes
func (e *EthereumHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all nodes")
}

// Create creates a single Ethereum node from spec
func (e *EthereumHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a node")
}

// Delete deletes a single Ethereum node by name
func (e *EthereumHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a node")
}

// Update updates a single node by name from new spec delta
func (e *EthereumHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a node")
}

// Register registers all routes on the given router
func (e *EthereumHandler) Register(router fiber.Router) {
	router.Post("/", e.Create)
	router.Get("/", e.List)
	router.Get("/:name", e.Get)
	router.Put("/:name", e.Update)
	router.Delete("/:name", e.Delete)
}

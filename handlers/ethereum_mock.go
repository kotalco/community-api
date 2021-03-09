package handlers

import "github.com/gofiber/fiber/v2"

// EthereumMockHandler is Ethereum Mock handler
type EthereumMockHandler struct{}

// NewEthereumMockHandler returns new mock handler
func NewEthereumMockHandler() Handler {
	return &EthereumMockHandler{}
}

// Get gets a single node
func (e *EthereumMockHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get mock node")
}

// List lists all Ethereum nodes
func (e *EthereumMockHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all mock nodes")
}

// Create creates a single Ethereum node from spec
func (e *EthereumMockHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create mock node")
}

// Delete deletes a single Ethereum node by name
func (e *EthereumMockHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete mock node")
}

// Update updates a single node by name from new spec delta
func (e *EthereumMockHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update mock node")
}

// Register registers all routes on the given router
func (e *EthereumMockHandler) Register(router fiber.Router) {
	router.Post("/", e.Create)
	router.Get("/", e.List)
	router.Get("/:name", e.Get)
	router.Put("/:name", e.Update)
	router.Delete("/:name", e.Delete)
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	models "github.com/kotalco/api/models/ethereum"
)

// NodesStore is in-memory nodes store
var NodesStore = map[string]*models.Node{}

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
	node := new(models.Node)

	if err := c.BodyParser(node); err != nil {
		return err
	}

	// check if node exist with this name
	if NodesStore[node.Name] != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(map[string]string{
			"error": fmt.Sprintf("node by name %s already exist", node.Name),
		})
	}

	NodesStore[node.Name] = node

	return c.JSON(node)
}

// Delete deletes a single Ethereum node by name
func (e *NodeMockHandler) Delete(c *fiber.Ctx) error {

	name := c.Params("name")

	// check if node exist with this name doesn't exist
	if NodesStore[name] == nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(map[string]string{
			"error": fmt.Sprintf("node by name %s doesn't exist", name),
		})
	}

	// remove node from the store
	delete(NodesStore, name)

	return c.SendStatus(http.StatusNoContent)
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

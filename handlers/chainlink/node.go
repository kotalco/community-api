package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/api/handlers"
	sharedHandlers "github.com/kotalco/api/handlers/shared"
	"github.com/kotalco/api/k8s"
	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
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
	name := c.Params("name")
	node := &chainlinkv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, node); err != nil {

		log.Print(err)

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("node by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get node by name %s", name),
		})
	}

	c.Locals("node", node)
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

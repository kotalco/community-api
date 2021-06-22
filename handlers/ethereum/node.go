package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/ethereum"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// NodeHandler is Ethereum node handler
type NodeHandler struct{}

// NewNodeHandler returns new Ethereum node handler
func NewNodeHandler() handlers.Handler {
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
	model := new(models.Node)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	node := &ethereumv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: "default",
		},
		Spec: ethereumv1alpha1.NodeSpec{
			NetworkConfig: ethereumv1alpha1.NetworkConfig{
				Join: model.Network,
			},
			Client: ethereumv1alpha1.EthereumClient(model.Client),
		},
	}

	var rpcAPI []ethereumv1alpha1.API
	if model.RPC != nil {
		rpc := *model.RPC
		if rpc {
			rpcAPI = []ethereumv1alpha1.API{}
			for _, api := range model.RPCAPI {
				rpcAPI = append(rpcAPI, ethereumv1alpha1.API(api))
			}
			node.Spec.RPCAPI = rpcAPI
		}
		node.Spec.RPC = rpc
	}

	var wsAPI []ethereumv1alpha1.API
	if model.WS != nil {
		ws := *model.WS
		if ws {
			wsAPI = []ethereumv1alpha1.API{}
			for _, api := range model.WSAPI {
				wsAPI = append(wsAPI, ethereumv1alpha1.API(api))
			}
			node.Spec.WSAPI = wsAPI
		}
		node.Spec.WS = ws
	}

	if err := k8s.Client().Create(c.Context(), node); err != nil {
		if errors.IsAlreadyExists(err) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": fmt.Sprintf("node by name %s already exist", model.Name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create node",
		})
	}

	return c.Status(http.StatusCreated).JSON(map[string]interface{}{
		"node": model,
	})
}

// Delete deletes a single Ethereum node by name
func (e *NodeHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a node")
}

// Update updates a single node by name from new spec delta
func (e *NodeHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a node")
}

// validateNodeExist validate node by name exist
func validateNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")
	node := &ethereumv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, node); err != nil {

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("node by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get node by name %s", name),
		})
	}

	return c.Next()

}

// Register registers all routes on the given router
func (e *NodeHandler) Register(router fiber.Router) {
	router.Post("/", e.Create)
	router.Get("/", e.List)
	router.Get("/:name", validateNodeExist, e.Get)
	router.Put("/:name", validateNodeExist, e.Update)
	router.Delete("/:name", validateNodeExist, e.Delete)
}

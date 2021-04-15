package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	models "github.com/kotalco/api/models/ethereum"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// nodesStore is in-memory nodes store
var nodesStore = map[string]*ethereumv1alpha1.Node{}

// NodeMockHandler is Ethereum node mock handler
type NodeMockHandler struct{}

// NewNodeMockHandler returns new Ethereum node mock handler
func NewNodeMockHandler() handlers.Handler {
	return &NodeMockHandler{}
}

// Get gets a single node
func (e *NodeMockHandler) Get(c *fiber.Ctx) error {
	name := c.Params("name")

	node := nodesStore[name]
	model := models.FromEthereumNode(node)

	return c.JSON(map[string]interface{}{
		"node": model,
	})
}

// List lists all Ethereum nodes
func (e *NodeMockHandler) List(c *fiber.Ctx) error {
	nodes := []models.Node{}

	for _, node := range nodesStore {
		model := models.FromEthereumNode(node)
		nodes = append(nodes, *model)
	}

	return c.JSON(map[string]interface{}{
		"nodes": nodes,
	})

}

// Create creates a single Ethereum node from spec
func (e *NodeMockHandler) Create(c *fiber.Ctx) error {
	model := new(models.Node)

	if err := c.BodyParser(model); err != nil {
		return err
	}

	// check if node exist with this name
	if nodesStore[model.Name] != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(map[string]string{
			"error": fmt.Sprintf("node by name %s already exist", model.Name),
		})
	}

	var rpcAPI []ethereumv1alpha1.API
	if model.RPC {
		rpcAPI = []ethereumv1alpha1.API{}
		for _, api := range model.RPCAPI {
			rpcAPI = append(rpcAPI, ethereumv1alpha1.API(api))
		}
	}

	var wsAPI []ethereumv1alpha1.API
	if model.WS {
		wsAPI = []ethereumv1alpha1.API{}
		for _, api := range model.WSAPI {
			wsAPI = append(wsAPI, ethereumv1alpha1.API(api))
		}
	}

	node := &ethereumv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: model.Name,
		},
		Spec: ethereumv1alpha1.NodeSpec{
			NetworkConfig: ethereumv1alpha1.NetworkConfig{
				Join: model.Network,
			},
			Client: ethereumv1alpha1.EthereumClient(model.Client),
			RPC:    model.RPC,
			RPCAPI: rpcAPI,
			WS:     model.WS,
			WSAPI:  wsAPI,
		},
	}

	// TODO: default the node

	nodesStore[model.Name] = node

	return c.Status(http.StatusCreated).JSON(map[string]interface{}{
		"node": model,
	})
}

// Delete deletes a single Ethereum node by name
func (e *NodeMockHandler) Delete(c *fiber.Ctx) error {

	name := c.Params("name")

	// remove node from the store
	delete(nodesStore, name)

	return c.SendStatus(http.StatusNoContent)
}

// Update updates a single node by name from new spec delta
func (e *NodeMockHandler) Update(c *fiber.Ctx) error {

	name := c.Params("name")

	model := new(models.Node)

	if err := c.BodyParser(model); err != nil {
		return err
	}

	if model.Client != "" {
		nodesStore[name].Spec.Client = ethereumv1alpha1.EthereumClient(model.Client)
	}

	if len(model.RPCAPI) != 0 {

		rpcAPI := []ethereumv1alpha1.API{}
		for _, api := range model.RPCAPI {
			rpcAPI = append(rpcAPI, ethereumv1alpha1.API(api))
		}
		nodesStore[name].Spec.RPCAPI = rpcAPI
	}

	nodesStore[name].Spec.RPC = model.RPC

	if len(model.WSAPI) != 0 {
		wsAPI := []ethereumv1alpha1.API{}
		for _, api := range model.RPCAPI {
			wsAPI = append(wsAPI, ethereumv1alpha1.API(api))
		}
		nodesStore[name].Spec.WSAPI = wsAPI
	}

	nodesStore[name].Spec.WS = model.WS

	updatedModel := models.FromEthereumNode(nodesStore[name])

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": updatedModel,
	})
}

// validateNodeExist validate node by name exist
func validateNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")

	if nodesStore[name] != nil {
		return c.Next()
	}
	return c.Status(http.StatusNotFound).JSON(map[string]string{
		"error": fmt.Sprintf("node by name %s doesn't exist", c.Params("name")),
	})
}

// Register registers all routes on the given router
func (e *NodeMockHandler) Register(router fiber.Router) {
	router.Post("/", e.Create)
	router.Get("/", e.List)
	router.Get("/:name", validateNodeExist, e.Get)
	router.Put("/:name", validateNodeExist, e.Update)
	router.Delete("/:name", validateNodeExist, e.Delete)
}

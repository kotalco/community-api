package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/ethereum"
	"github.com/kotalco/api/shared"
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
	node := c.Locals("node").(*ethereumv1alpha1.Node)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": models.FromEthereumNode(node),
	})
}

// List lists all Ethereum nodes
func (e *NodeHandler) List(c *fiber.Ctx) error {
	nodes := &ethereumv1alpha1.NodeList{}
	if err := k8s.Client().List(c.Context(), nodes); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all nodes",
		})
	}

	nodeModels := []models.Node{}

	page := c.Query("page")
	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	start, end := shared.Page(uint(len(nodes.Items)), uint(p))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[i].CreationTimestamp.Before(&nodes.Items[j].CreationTimestamp)
	})

	for _, node := range nodes.Items[start:end] {
		nodeModels = append(nodeModels, *models.FromEthereumNode(&node))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"nodes": nodeModels,
	})

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

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Create(c.Context(), node); err != nil {
		log.Println(err)
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
		"node": models.FromEthereumNode(node),
	})
}

// Delete deletes a single Ethereum node by name
func (e *NodeHandler) Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereumv1alpha1.Node)

	if err := k8s.Client().Delete(c.Context(), node); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete node by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates a single node by name from new spec delta
func (e *NodeHandler) Update(c *fiber.Ctx) error {
	model := new(models.Node)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	name := c.Params("name")
	node := c.Locals("node").(*ethereumv1alpha1.Node)

	if model.SyncMode != "" {
		node.Spec.SyncMode = ethereumv1alpha1.SynchronizationMode(model.SyncMode)
	}
	if model.P2PPort != 0 {
		node.Spec.P2PPort = model.P2PPort
	}
	if model.RPC != nil {
		node.Spec.RPC = *model.RPC
	}
	if node.Spec.RPC {
		if len(model.RPCAPI) != 0 {
			rpcAPI := []ethereumv1alpha1.API{}
			for _, api := range model.RPCAPI {
				rpcAPI = append(rpcAPI, ethereumv1alpha1.API(api))
			}
			node.Spec.RPCAPI = rpcAPI
		}
		if model.RPCPort != 0 {
			node.Spec.RPCPort = model.RPCPort
		}
	}

	if model.WS != nil {
		node.Spec.WS = *model.WS
	}
	if node.Spec.WS {
		if len(model.WSAPI) != 0 {
			wsAPI := []ethereumv1alpha1.API{}
			for _, api := range model.WSAPI {
				wsAPI = append(wsAPI, ethereumv1alpha1.API(api))
			}
			node.Spec.WSAPI = wsAPI
		}
		if model.WSPort != 0 {
			node.Spec.WSPort = model.WSPort
		}
	}

	if model.GraphQL != nil {
		node.Spec.GraphQL = *model.GraphQL
	}
	if node.Spec.GraphQL {
		if model.GraphQLPort != 0 {
			node.Spec.GraphQLPort = model.GraphQLPort
		}
	}

	if len(model.Hosts) != 0 {
		node.Spec.Hosts = model.Hosts
	}

	if len(model.CORSDomains) != 0 {
		node.Spec.CORSDomains = model.CORSDomains
	}

	if model.CPU != "" {
		node.Spec.CPU = model.CPU
	}
	if model.CPULimit != "" {
		node.Spec.CPULimit = model.CPULimit
	}
	if model.Memory != "" {
		node.Spec.Memory = model.Memory
	}
	if model.MemoryLimit != "" {
		node.Spec.MemoryLimit = model.MemoryLimit
	}
	if model.Storage != "" {
		node.Spec.Storage = model.Storage
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Update(c.Context(), node); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't update node by name %s", name),
		})
	}

	updatedModel := models.FromEthereumNode(node)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": updatedModel,
	})
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

	c.Locals("node", node)

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

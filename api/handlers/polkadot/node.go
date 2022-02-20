package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kotalco/api/api/handlers"
	shared2 "github.com/kotalco/api/api/handlers/shared"
	"github.com/kotalco/api/internal/models/polkadot"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/shared"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	"github.com/ybbus/jsonrpc/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NodeHandler is Polkadot node handler
type NodeHandler struct{}

// NewNodeHandler creates a new Polkadot node handler
func NewNodeHandler() handlers.Handler {
	return &NodeHandler{}
}

// Get gets a single Polkadot node by name
func (n *NodeHandler) Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*polkadotv1alpha1.Node)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": models.FromPolkadotNode(node),
	})
}

// List returns all Polkadot nodes
func (n *NodeHandler) List(c *fiber.Ctx) error {
	nodes := &polkadotv1alpha1.NodeList{}
	if err := k8s.Client().List(c.Context(), nodes, client.InNamespace("default")); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all nodes",
		})
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	nodeModels := []models.Node{}

	page, _ := strconv.Atoi(c.Query("page"))

	start, end := shared.Page(uint(len(nodes.Items)), uint(page))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	for _, node := range nodes.Items[start:end] {
		nodeModels = append(nodeModels, *models.FromPolkadotNode(&node))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"nodes": nodeModels,
	})

}

// Create creates Polkadot node from spec
func (n *NodeHandler) Create(c *fiber.Ctx) error {
	model := new(models.Node)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	node := &polkadotv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: "default",
		},
		Spec: polkadotv1alpha1.NodeSpec{
			Network: model.Network,
			RPC:     true,
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

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"node": models.FromPolkadotNode(node),
	})
}

// Delete deletes Polkadot node by name
func (n *NodeHandler) Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*polkadotv1alpha1.Node)

	if err := k8s.Client().Delete(c.Context(), node); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete node by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates Polkadot node by name from spec
func (n *NodeHandler) Update(c *fiber.Ctx) error {
	model := new(models.Node)
	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	name := c.Params("name")
	node := c.Locals("node").(*polkadotv1alpha1.Node)

	if model.NodePrivateKeySecretName != "" {
		node.Spec.NodePrivateKeySecretName = model.NodePrivateKeySecretName
	}

	if model.Validator != nil {
		node.Spec.Validator = *model.Validator
	}

	if model.SyncMode != "" {
		node.Spec.SyncMode = polkadotv1alpha1.SynchronizationMode(model.SyncMode)
	}

	if model.Pruning != nil {
		node.Spec.Pruning = model.Pruning
	}

	if model.P2PPort != 0 {
		node.Spec.P2PPort = model.P2PPort
	}

	if model.RetainedBlocks != 0 {
		node.Spec.RetainedBlocks = model.RetainedBlocks
	}

	if model.Logging != "" {
		node.Spec.Logging = sharedAPI.VerbosityLevel(model.Logging)
	}

	if model.Telemetry != nil {
		node.Spec.Telemetry = *model.Telemetry
	}

	if model.TelemetryURL != "" {
		node.Spec.TelemetryURL = model.TelemetryURL
	}

	if model.Prometheus != nil {
		node.Spec.Prometheus = *model.Prometheus
	}

	if model.PrometheusPort != 0 {
		node.Spec.PrometheusPort = model.PrometheusPort
	}

	if model.RPC != nil {
		node.Spec.RPC = *model.RPC
	}

	if model.RPCPort != 0 {
		node.Spec.RPCPort = model.RPCPort
	}

	if model.WS != nil {
		node.Spec.WS = *model.WS
	}

	if model.WSPort != 0 {
		node.Spec.WSPort = model.WSPort
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

	updatedModel := models.FromPolkadotNode(node)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": updatedModel,
	})
}

// Count returns total number of nodes
func (n *NodeHandler) Count(c *fiber.Ctx) error {
	nodes := &polkadotv1alpha1.NodeList{}
	if err := k8s.Client().List(c.Context(), nodes, client.InNamespace("default")); err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	return c.SendStatus(http.StatusOK)
}

func (e *NodeHandler) Stats(c *websocket.Conn) {
	defer c.Close()

	type Result struct {
		Error string `json:"error,omitempty"`
		// system_syncState call
		CurrentBlock uint `json:"currentBlock,omitempty"`
		HighestBlock uint `json:"highestBlock,omitempty"`
		// system_health call
		Peers   uint `json:"peersCount,omitempty"`
		Syncing bool `json:"syncing"`
	}

	// Mock serever
	if os.Getenv("MOCK") == "true" {
		var currentBlock, highestBlock, peersCount uint
		for {
			currentBlock += 3
			highestBlock += 32
			peersCount += 1

			r := &Result{
				CurrentBlock: currentBlock,
				HighestBlock: highestBlock,
				Peers:        peersCount,
				Syncing:      peersCount%4 != 0,
			}

			var msg []byte

			if peersCount > 40 {
				peersCount = 0
				r = &Result{
					Error: "JSON-RPC server is not enabled",
				}
			}

			msg, _ = json.Marshal(r)
			c.WriteMessage(websocket.TextMessage, []byte(msg))
			time.Sleep(time.Second)
		}
	}

	name := c.Params("name")
	node := &polkadotv1alpha1.Node{}
	key := types.NamespacedName{
		Namespace: "default",
		Name:      name,
	}

	for {

		err := k8s.Client().Get(context.Background(), key, node)
		if errors.IsNotFound(err) {
			c.WriteJSON(fiber.Map{
				"error": fmt.Sprintf("node by name %s doesn't exist", name),
			})
			return
		}

		if !node.Spec.RPC {
			c.WriteJSON(fiber.Map{
				"error": "JSON-RPC server is not enabled",
			})
			time.Sleep(time.Second)
			continue
		}

		client := jsonrpc.NewClient(fmt.Sprintf("http://%s:%d", node.Name, node.Spec.RPCPort))

		type SyncState struct {
			CurrentBlock uint `json:"currentBlock"`
			HighestBlock uint `json:"highestBlock"`
		}

		// sync state rpc call
		syncState := &SyncState{}
		err = client.CallFor(syncState, "system_syncState")
		if err != nil {
			fmt.Println(err)
		}

		// system_health .isSyncing .peers
		type SystemHealth struct {
			Syncing    bool `json:"isSyncing"`
			PeersCount uint `json:"peers"`
		}

		// system health rpc call
		systemHealth := &SystemHealth{}
		err = client.CallFor(systemHealth, "system_health")
		if err != nil {
			fmt.Println(err)
		}

		c.WriteJSON(fiber.Map{
			"currentBlock": syncState.CurrentBlock,
			"highestBlock": syncState.HighestBlock,
			"peersCount":   systemHealth.PeersCount,
			"syncing":      systemHealth.Syncing,
		})

		time.Sleep(time.Second)
	}
}

// validateNodeExist validates Polkadot node by name exist
func validateNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")
	node := &polkadotv1alpha1.Node{}
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

// Register registers all handlers on the given router
func (n *NodeHandler) Register(router fiber.Router) {
	router.Post("/", n.Create)
	router.Head("/", n.Count)
	router.Get("/", n.List)
	router.Get("/:name", validateNodeExist, n.Get)
	router.Get("/:name/logs", websocket.New(shared2.Logger))
	router.Get("/:name/status", websocket.New(shared2.Status))
	router.Get("/:name/stats", websocket.New(n.Stats))
	router.Put("/:name", validateNodeExist, n.Update)
	router.Delete("/:name", validateNodeExist, n.Delete)
}

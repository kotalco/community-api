package polkadot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/internal/polkadot"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	polkadotv1alpha1 "github.com/kotalco/kotal/apis/polkadot/v1alpha1"
	"github.com/ybbus/jsonrpc/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

const (
	nameKeyword = "name"
)

var (
	k8sClient = k8s.NewClientService()
	service   = polkadot.NewPolkadotService()
)

// Get gets a single Polkadot node by name
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*polkadotv1alpha1.Node)

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(polkadot.PolkadotDto).FromPolkadotNode(node)))
}

// List returns all Polkadot nodes
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	nodes, err := service.List(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	start, end := shared.Page(uint(len(nodes.Items)), uint(page), uint(limit))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(polkadot.PolkadotListDto).FromPolkadotNode(nodes.Items[start:end])))
}

// Create creates Polkadot node from spec
func Create(c *fiber.Ctx) error {
	dto := new(polkadot.PolkadotDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	dto.Namespace = c.Locals("namespace").(string)

	err := dto.MetaDataDto.Validate()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	node, err := service.Create(dto)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(polkadot.PolkadotDto).FromPolkadotNode(node)))
}

// Delete deletes Polkadot node by name
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*polkadotv1alpha1.Node)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates Polkadot node by name from spec
func Update(c *fiber.Ctx) error {
	dto := new(polkadot.PolkadotDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(err)
	}

	node := c.Locals("node").(*polkadotv1alpha1.Node)

	node, err := service.Update(dto, node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(polkadot.PolkadotDto).FromPolkadotNode(node)))
}

// Count returns total number of nodes
func Count(c *fiber.Ctx) error {
	length, err := service.Count(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", *length))

	return c.SendStatus(http.StatusOK)
}

func Stats(c *websocket.Conn) {
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
	ns := c.Locals("namespace").(string)
	nodeKey := types.NamespacedName{
		Namespace: ns,
		Name:      name,
	}
	podKey := types.NamespacedName{
		Namespace: ns,
		Name:      fmt.Sprintf("%s-0", name),
	}

nodeCheck:
	// checking node is found and rpc is enabled
	// because node can be deleted, and rpc closed
	// during the lifetime of socket connection
	node := &polkadotv1alpha1.Node{}
	if err := k8sClient.Get(context.Background(), nodeKey, node); errors.IsNotFound(err) {
		c.WriteJSON(fiber.Map{
			"error": fmt.Sprintf("node by name %s doesn't exist", name),
		})
		return
	}

	if !node.Spec.RPC {
		if err := c.WriteJSON(fiber.Map{"error": "JSON-RPC server is not enabled"}); err != nil {
			return
		}
		time.Sleep(3 * time.Second)
		goto nodeCheck
	}

podCheck:
	// check pod exist if any rpc failed
	// if pod is not found, check if node has been deleted
	pod := &corev1.Pod{}
	if err := k8sClient.Get(context.Background(), podKey, pod); err != nil {
		if apierrors.IsNotFound(err) {
			goto nodeCheck
		}
	}
	if pod.Status.Phase != corev1.PodRunning {
		time.Sleep(3 * time.Second)
		goto podCheck
	}

	endpoint := fmt.Sprintf("http://%s.%s:%d", node.Name, node.Namespace, node.Spec.RPCPort)
	rpcClient := jsonrpc.NewClient(endpoint)

	for {

		type SyncState struct {
			CurrentBlock uint `json:"currentBlock"`
			HighestBlock uint `json:"highestBlock"`
		}

		// sync state rpc call
		syncState := &SyncState{}
		if err := rpcClient.CallFor(syncState, "system_syncState"); err != nil {
			time.Sleep(3 * time.Second)
			goto podCheck
		}

		// system_health .isSyncing .peers
		type SystemHealth struct {
			Syncing    bool `json:"isSyncing"`
			PeersCount uint `json:"peers"`
		}

		// system health rpc call
		systemHealth := &SystemHealth{}
		if err := rpcClient.CallFor(systemHealth, "system_health"); err != nil {
			time.Sleep(3 * time.Second)
			goto podCheck
		}

		if err := c.WriteJSON(fiber.Map{
			"currentBlock": syncState.CurrentBlock,
			"highestBlock": syncState.HighestBlock,
			"peersCount":   systemHealth.PeersCount,
			"syncing":      systemHealth.Syncing,
		}); err != nil {
			return
		}

		time.Sleep(time.Second)
	}
}

// ValidateNodeExist validates Polkadot node by name exist
func ValidateNodeExist(c *fiber.Ctx) error {
	nameSpacedName := types.NamespacedName{
		Name:      c.Params(nameKeyword),
		Namespace: c.Locals("namespace").(string),
	}

	node, err := service.Get(nameSpacedName)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("node", node)

	return c.Next()
}

package near

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/api/internal/near"
	restErrors "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/shared"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
	"github.com/ybbus/jsonrpc/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	NODE_NAME_KEYWORD = "name"
	NAMESPACE_KEYWORD = "namespace"
	DEFAULT_NAMESPACE = "default"
)

var service = near.NearService

// Get gets a single NEAR node by name
// 1-get the node validated from ValidateNodeExist method
// 2-marshall node to dto and format the response
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*nearv1alpha1.Node)

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(near.NearDto).FromNEARNode(node)))
}

// List returns all NEAR nodes
// 1-get the pagination qs default to 0
// 2-call service to return node models
// 3-make the pagination
// 3-marshall nodes  to near dto and format the response using NewResponse
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))

	nodes, err := service.List(c.Query(NAMESPACE_KEYWORD, DEFAULT_NAMESPACE))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	start, end := shared.Page(uint(len(nodes.Items)), uint(page))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(near.NearListDto).FromNEARNode(nodes.Items[start:end])))
}

// Create creates NEAR node from spec
// Create creates near node from spec
// 1-Todo validate request body and return validation error
// 2-call near service to create near node
// 2-marshall node to dto and format the response
func Create(c *fiber.Ctx) error {
	dto := new(near.NearDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	node, err := service.Create(dto)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(near.NearDto).FromNEARNode(node)))
}

// Delete deletes NEAR node by name
// 1-get node from locals which checked and assigned by ValidateNodeExist
// 2-call near service to delete the node
// 3-return ok if deleted with no errors
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*nearv1alpha1.Node)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates NEAR node by name from spec
// 1-todo validate request body and return validation errors if exits
// 2-get node from locals which checked and assigned by ValidateNodeExist
// 3-call near service to update node which returns *nearv1alpha1.Node
// 4-marshall node to node dto and format the response
func Update(c *fiber.Ctx) error {
	dto := new(near.NearDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	node := c.Locals("node").(*nearv1alpha1.Node)

	node, err := service.Update(dto, node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(near.NearDto).FromNEARNode(node)))
}

// Count returns total number of nodes
// 1-call near service to get exiting node list
// 2-create X-Total-Count header with the length
// 3-return
func Count(c *fiber.Ctx) error {
	length, err := service.Count(c.Query(NAMESPACE_KEYWORD, DEFAULT_NAMESPACE))
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
		// network_info call
		ActivePeersCount       uint `json:"activePeersCount,omitempty"`
		MaxPeersCount          uint `json:"maxPeersCount,omitempty"`
		SentBytesPerSecond     uint `json:"sentBytesPerSecond,omitempty"`
		ReceivedBytesPerSecond uint `json:"receivedBytesPerSecond,omitempty"`
		// status call
		LatestBlockHeight   uint `json:"latestBlockHeight,omitempty"`
		EarliestBlockHeight uint `json:"earliestBlockHeight,omitempty"`
		Syncing             bool `json:"syncing,omitempty"`
	}

	// Mock serever
	if os.Getenv("MOCK") == "true" {
		var activePeersCount, sentBytesPerSecond, receivedBytesPerSecond, latestBlockHeight, earliestBlockHeight uint
		for {
			activePeersCount++
			sentBytesPerSecond += 100
			receivedBytesPerSecond += 100
			latestBlockHeight += 36
			earliestBlockHeight += 3

			r := &Result{
				ActivePeersCount:       activePeersCount,
				MaxPeersCount:          40,
				SentBytesPerSecond:     sentBytesPerSecond,
				ReceivedBytesPerSecond: receivedBytesPerSecond,
				LatestBlockHeight:      latestBlockHeight,
				EarliestBlockHeight:    earliestBlockHeight,
				Syncing:                true,
			}

			var msg []byte

			if activePeersCount > 40 {
				activePeersCount = 10
				r = &Result{
					Error: "rpc is not enabled",
				}
			}

			msg, _ = json.Marshal(r)
			c.WriteMessage(websocket.TextMessage, []byte(msg))
			time.Sleep(time.Second)
		}
	}

	name := c.Params("name")
	node := &nearv1alpha1.Node{}
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

		type NodeStatus struct {
			SyncInfo struct {
				LatestBlockHeight   uint `json:"latest_block_height"`
				EarliestBlockHeight uint `json:"earliest_block_height"`
				Syncing             bool `json:"syncing"`
			} `json:"sync_info"`
		}

		// node status rpc call
		nodeStatus := &NodeStatus{}
		err = client.CallFor(nodeStatus, "status")
		if err != nil {
			fmt.Println(err)
		}

		type NetworkInfo struct {
			ActivePeersCount       uint `json:"num_active_peers"`
			MaxPeersCount          uint `json:"peer_max_count"`
			SentBytesPerSecond     uint `json:"sent_bytes_per_sec"`
			ReceivedBytesPerSecond uint `json:"received_bytes_per_sec"`
		}

		// network info rpc call
		networkInfo := &NetworkInfo{}
		err = client.CallFor(networkInfo, "network_info")
		if err != nil {
			fmt.Println(err)
		}

		c.WriteJSON(fiber.Map{
			"activePeersCount":       networkInfo.ActivePeersCount,
			"maxPeersCount":          networkInfo.MaxPeersCount,
			"sentBytesPerSecond":     networkInfo.SentBytesPerSecond,
			"receivedBytesPerSecond": networkInfo.ReceivedBytesPerSecond,
			"latestBlockHeight":      nodeStatus.SyncInfo.LatestBlockHeight,
			"earliestBlockHeight":    nodeStatus.SyncInfo.EarliestBlockHeight,
			"syncing":                nodeStatus.SyncInfo.Syncing,
		})

		time.Sleep(time.Second)
	}
}

// ValidateNodeExist  validate node by name exist acts as a validation for all handlers the needs to find near node by name
// 1-call near service to check if node exits
// 2-return Not found if it's not
// 3-save the node to local with the key node to be used by the other handlers
func ValidateNodeExist(c *fiber.Ctx) error {
	nameSpacedName := types.NamespacedName{
		Name:      c.Params(NODE_NAME_KEYWORD),
		Namespace: c.Query(NAMESPACE_KEYWORD, DEFAULT_NAMESPACE),
	}

	node, err := service.Get(nameSpacedName)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("node", node)

	return c.Next()
}

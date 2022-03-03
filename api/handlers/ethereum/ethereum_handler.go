// Package ethereum  handler is the representation layer for the  ethereum node
// uses the k8 client to CRUD the nodes
package ethereum

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/api/internal/ethereum"
	restErrors "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/shared"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	"github.com/ybbus/jsonrpc/v2"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var service = ethereum.EthereumService

// Get returns a single ethereum node by name
// 1-get the node validated from ValidateNodeExist method
// 2-marshall node to dto and format the response
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereumv1alpha1.Node)

	return c.JSON(shared.NewResponse(new(ethereum.EthereumDto).FromEthereumNode(node)))
}

// Create creates ethereum node from the given spec
// 1-Todo validate request body and return validation error
// 2-call chain link service to create ethereum node
// 2-marshall node to dto and format the response
func Create(c *fiber.Ctx) error {
	dto := new(ethereum.EthereumDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	node, err := service.Create(dto)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(ethereum.EthereumDto).FromEthereumNode(node)))
}

// Update updates a single ethereum node by name from spec
// 1-todo validate request body and return validation errors if exits
// 2-get node from locals which checked and assigned by ValidateNodeExist
// 3-call ethereum service to update node which returns *ethereumv1alpha1.Node
// 4-marshall node to node dto and format the response
func Update(c *fiber.Ctx) error {
	dto := new(ethereum.EthereumDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	node := c.Locals("node").(*ethereumv1alpha1.Node)

	node, err := service.Update(dto, node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(ethereum.EthereumDto).FromEthereumNode(node)))
}

// List returns all ethereum nodes
// 1-get the pagination qs default to 0
// 2-call service to return node models
// 3-make the pagination
// 3-marshall nodes  to ethereum dto and format the response using NewResponse
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))

	nodes, err := service.List()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	start, end := shared.Page(uint(len(nodes.Items)), uint(page))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(ethereum.EthereumListDto).FromEthereumNode(nodes.Items[start:end])))
}

// Delete a single ethereum node by name
// 1-get node from locals which checked and assigned by ValidateNodeExist
// 2-call ethereum service to delete the node
// 3-return ok if deleted with no errors
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereumv1alpha1.Node)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Count returns total number of nodes
// 1-call ethereum service to get exiting node list
// 2-create X-Total-Count header with the length
// 3-return
func Count(c *fiber.Ctx) error {
	length, err := service.Count()
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
		// eth_syncing call
		CurrentBlock uint `json:"currentBlock,omitempty"`
		HighestBlock uint `json:"highestBlock,omitempty"`
		// net_peerCount call
		Peers uint `json:"peersCount,omitempty"`
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
			}

			var msg []byte

			if peersCount > 20 {
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

	for {

		node, err := service.Get(name)

		if err != nil {
			c.WriteJSON(err)
			return
		}

		if !node.Spec.RPC {
			c.WriteJSON(restErrors.NewBadRequestError("rpc is not enabled"))
			return
		}

		client := jsonrpc.NewClient(fmt.Sprintf("http://%s:%d", node.Name, node.Spec.RPCPort))

		type SyncStatus struct {
			CurrentBlock string `json:"currentBlock"`
			HighestBlock string `json:"highestBlock"`
		}

		// sync status
		syncStatus := SyncStatus{}
		client.CallFor(&syncStatus, "eth_syncing")

		current := new(big.Int)
		current.SetString(strings.Replace(syncStatus.CurrentBlock, "0x", "", 1), 16)

		highest := new(big.Int)
		highest.SetString(strings.Replace(syncStatus.HighestBlock, "0x", "", 1), 16)

		// peer count
		var peerCount string
		client.CallFor(&peerCount, "net_peerCount")

		count := new(big.Int)
		count.SetString(strings.Replace(peerCount, "0x", "", 1), 16)

		c.WriteJSON(fiber.Map{
			"currentBlock": current.String(),
			"highestBlock": highest.String(),
			"peersCount":   count,
		})

		time.Sleep(time.Second)
	}
}

// ValidateNodeExist validate node by name exist acts as a validation for all handlers the needs to find ethereum by name
// 1-call ethereum service to check if node exits
// 2-return 404 if it's not
// 3-save the node to local with the key node to be used by the other handlers
func ValidateNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")

	node, err := service.Get(name)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("node", node)
	return c.Next()
}

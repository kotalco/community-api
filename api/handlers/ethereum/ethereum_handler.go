// Package ethereum  handler is the representation layer for the  ethereum node
// uses the k8 client to CRUD the nodes
package ethereum

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/api/internal/ethereum"
	restError "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/shared"
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	"github.com/ybbus/jsonrpc/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var service = ethereum.EthereumService

// Get returns a single chainlink node by name
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereumv1alpha1.Node)

	return c.JSON(shared.NewResponse(new(Node).FromEthereumNode(node)))
}

// Create creates chainlink node from the given spec
func Create(c *fiber.Ctx) error {
	request := new(Node)
	if err := c.BodyParser(request); err != nil {
		//Todo add Validation
		badReq := restError.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	node := &ethereumv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.Name,
			Namespace: "default",
		},
		Spec: ethereumv1alpha1.NodeSpec{
			Network:                  request.Network,
			Client:                   ethereumv1alpha1.EthereumClient(request.Client),
			RPC:                      true,
			NodePrivateKeySecretName: request.NodePrivateKeySecretName,
			Resources: sharedAPI.Resources{
				StorageClass: request.StorageClass,
			},
		},
	}

	node, err := service.Create(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(Node).FromEthereumNode(node)))
}

// Update updates a single chainlink node by name from spec
func Update(c *fiber.Ctx) error {
	request := new(Node)
	if err := c.BodyParser(request); err != nil {
		//Todo add Validation
		badReq := restError.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	node := c.Locals("node").(*ethereumv1alpha1.Node)
	if request.Logging != "" {
		node.Spec.Logging = sharedAPI.VerbosityLevel(request.Logging)
	}
	if request.NodePrivateKeySecretName != "" {
		node.Spec.NodePrivateKeySecretName = request.NodePrivateKeySecretName
	}
	if request.SyncMode != "" {
		node.Spec.SyncMode = ethereumv1alpha1.SynchronizationMode(request.SyncMode)
	}
	if request.P2PPort != 0 {
		node.Spec.P2PPort = request.P2PPort
	}

	if request.Miner != nil {
		node.Spec.Miner = *request.Miner
	}
	if node.Spec.Miner {
		if request.Coinbase != "" {
			node.Spec.Coinbase = ethereumv1alpha1.EthereumAddress(request.Coinbase)
		}
		if request.Import != nil {
			node.Spec.Import = &ethereumv1alpha1.ImportedAccount{
				PrivateKeySecretName: request.Import.PrivateKeySecretName,
				PasswordSecretName:   request.Import.PasswordSecretName,
			}
		}
	}

	if request.RPC != nil {
		node.Spec.RPC = *request.RPC
	}
	if node.Spec.RPC {
		if len(request.RPCAPI) != 0 {
			rpcAPI := []ethereumv1alpha1.API{}
			for _, api := range request.RPCAPI {
				rpcAPI = append(rpcAPI, ethereumv1alpha1.API(api))
			}
			node.Spec.RPCAPI = rpcAPI
		}
		if request.RPCPort != 0 {
			node.Spec.RPCPort = request.RPCPort
		}
	}

	if request.WS != nil {
		node.Spec.WS = *request.WS
	}
	if node.Spec.WS {
		if len(request.WSAPI) != 0 {
			wsAPI := []ethereumv1alpha1.API{}
			for _, api := range request.WSAPI {
				wsAPI = append(wsAPI, ethereumv1alpha1.API(api))
			}
			node.Spec.WSAPI = wsAPI
		}
		if request.WSPort != 0 {
			node.Spec.WSPort = request.WSPort
		}
	}

	if request.GraphQL != nil {
		node.Spec.GraphQL = *request.GraphQL
	}
	if node.Spec.GraphQL {
		if request.GraphQLPort != 0 {
			node.Spec.GraphQLPort = request.GraphQLPort
		}
	}

	if len(request.Hosts) != 0 {
		node.Spec.Hosts = request.Hosts
	}

	if len(request.CORSDomains) != 0 {
		node.Spec.CORSDomains = request.CORSDomains
	}

	var bootnodes, staticNodes []ethereumv1alpha1.Enode

	if request.Bootnodes != nil {
		for _, bootnode := range *request.Bootnodes {
			bootnodes = append(bootnodes, ethereumv1alpha1.Enode(bootnode))
		}
	}
	node.Spec.Bootnodes = bootnodes

	if request.StaticNodes != nil {
		for _, staticNode := range *request.StaticNodes {
			staticNodes = append(staticNodes, ethereumv1alpha1.Enode(staticNode))
		}
	}
	node.Spec.StaticNodes = staticNodes

	if request.CPU != "" {
		node.Spec.CPU = request.CPU
	}
	if request.CPULimit != "" {
		node.Spec.CPULimit = request.CPULimit
	}
	if request.Memory != "" {
		node.Spec.Memory = request.Memory
	}
	if request.MemoryLimit != "" {
		node.Spec.MemoryLimit = request.MemoryLimit
	}
	if request.Storage != "" {
		node.Spec.Storage = request.Storage
	}

	node, err := service.Update(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(Node).FromEthereumNode(node)))
}

// List returns all chainlink nodes
func List(c *fiber.Ctx) error {
	nodes, err := service.List()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	// default page to 0
	//Todo move to the service =>create pagination package
	page, _ := strconv.Atoi(c.Query("page"))
	start, end := shared.Page(uint(len(nodes.Items)), uint(page))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	nodeModels := new(Nodes).FromEthereumNode(nodes.Items[start:end])
	return c.Status(http.StatusOK).JSON(shared.NewResponse(nodeModels))
}

// Delete a single chainlink node by name
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereumv1alpha1.Node)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Count returns total number of nodes
func Count(c *fiber.Ctx) error {
	nodes, err := service.List()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

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
			c.WriteJSON(restError.NewBadRequestError("rpc is not enabled"))
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

// validateNodeExist validate node by name exist
func ValidateNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")

	node, err := service.Get(name)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("node", node)
	return c.Next()
}

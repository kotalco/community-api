package beacon_node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/internal/ethereum2/beacon_node"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type request struct {
	url  string
	name string
}
type result struct {
	err  error
	name string
	data []byte
}

const (
	nameKeyword = "name"
)

var (
	service   = beacon_node.NewBeaconNodeService()
	k8sClient = k8s.NewClientService()
)

// Get gets a single ethereum 2.0 beacon node by name
// 1-get the node validated from ValidateNodeExist method
// 2-marshall node to dto and format the response
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereum2v1alpha1.BeaconNode)

	return c.JSON(shared.NewResponse(new(beacon_node.BeaconNodeDto).FromEthereum2BeaconNode(node)))
}

// List returns all ethereum 2.0 beacon nodes
// 1-get the pagination qs default to 0
// 2-call service to return node models
// 3-make the pagination
// 3-marshall nodes  to beacon node dto and format the response using NewResponse
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

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(beacon_node.BeaconNodeListDto).FromEthereum2BeaconNode(nodes.Items[start:end])))
}

// Create creates ethereum 2.0 beacon node from spec
// 1-Todo validate request body and return validation error
// 2-call beacon node service to create beacon node
// 2-marshall node to dto and format the response
func Create(c *fiber.Ctx) error {
	dto := new(beacon_node.BeaconNodeDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(err)
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

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(beacon_node.BeaconNodeDto).FromEthereum2BeaconNode(node)))
}

// Delete deletes ethereum 2.0 beacon node by name
// 1-get node from locals which checked and assigned by ValidateNodeExist
// 2-call beacon node service to delete the node
// 3-return ok if deleted with no errors
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereum2v1alpha1.BeaconNode)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates ethereum 2.0 beacon node by name from spec
// 1-todo validate request body and return validation errors if exits
// 2-get node from locals which checked and assigned by ValidateNodeExist
// 3-call beacon node  service to update node which returns *ethereum2v1alpha1.BeaconNode
// 4-marshall node to node dto and format the response
func Update(c *fiber.Ctx) error {
	dto := new(beacon_node.BeaconNodeDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid reqeust body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	beaconnode := c.Locals("node").(*ethereum2v1alpha1.BeaconNode)

	beaconnode, err := service.Update(dto, beaconnode)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(beacon_node.BeaconNodeDto).FromEthereum2BeaconNode(beaconnode)))
}

// Count returns total number of beacon nodes
// 1-call beacon node service to get exiting node list
// 2-create X-Total-Count header with the length
// 3-return
func Count(c *fiber.Ctx) error {
	length, err := service.Count(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", *length))

	return c.SendStatus(http.StatusOK)
}

// ValidateBeaconNodeExist  validate node by name exist acts as a validation for all handlers the needs to find beacon node by name
// 1-call beacon node service to check if node exits
// 2-return 404 if it's not
// 3-save the node to local with the key node to be used by the other handlers
func ValidateBeaconNodeExist(c *fiber.Ctx) error {
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

// Stats returns a websocket that emits peer  count and node syncing status
func Stats(c *websocket.Conn) {
	defer c.Close()

	name := c.Params("name")
	beaconnode := &ethereum2v1alpha1.BeaconNode{}
	nameSpacedName := types.NamespacedName{
		Namespace: c.Locals("namespace").(string),
		Name:      name,
	}
	err := k8sClient.Get(context.Background(), nameSpacedName, beaconnode)
	if err != nil {
		if errors.IsNotFound(err) {
			c.WriteJSON(fiber.Map{
				"error": fmt.Sprintf("peer by name %s doesn't exist", name),
			})
			return
		}
		c.WriteJSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	if !beaconnode.Spec.REST {
		c.WriteJSON(fiber.Map{
			"error": "node rest is not enabled",
		})
		return
	}

	for {
		jobs := make(chan request, 2)
		results := make(chan result, 2)

		for i := 0; i < 2; i++ {
			go worker(jobs, results)
		}

		baseUrl := fmt.Sprintf("http://%s:%d/eth/v1/node/", "localhost", beaconnode.Spec.RESTPort)
		jobs <- request{name: "peers", url: fmt.Sprintf("%speer_count", baseUrl)}
		jobs <- request{name: "isSyncing", url: fmt.Sprintf("%ssyncing", baseUrl)}

		close(jobs)

		type NodePeerCountDto struct {
			Disconnected  string `json:"disconnected"`
			Connecting    string `json:"connecting"`
			Connected     string `json:"connected"`
			Disconnecting string `json:"disconnecting"`
		}
		type NodeSyncingStatusDto struct {
			HeadSlot     string `json:"head_slot"`
			SyncDistance string `json:"sync_distance"`
			IsSyncing    bool   `json:"is_syncing"`
			IsOptimistic bool   `json:"is_optimistic"`
		}
		var nodeStatResponseDto struct {
			Peers   NodePeerCountDto     `json:"peers"`
			Syncing NodeSyncingStatusDto `json:"syncing"`
		}

		newNodeResponse := nodeStatResponseDto

		for i := 0; i < 2; i++ {
			resp := <-results
			if resp.err != nil {
				c.WriteJSON(fiber.Map{
					"error": err.Error(),
				})
				continue
			}
			switch resp.name {
			case "peers":
				var responseBody struct {
					Data NodePeerCountDto `json:"data"`
				}
				err = json.Unmarshal(resp.data, &responseBody)
				if err != nil {
					c.WriteJSON(fiber.Map{
						"error": err.Error(),
					})
					break
				}
				newNodeResponse.Peers.Disconnected = responseBody.Data.Disconnected
				newNodeResponse.Peers.Connecting = responseBody.Data.Connecting
				newNodeResponse.Peers.Connected = responseBody.Data.Connected
				newNodeResponse.Peers.Disconnecting = responseBody.Data.Disconnecting
				break
			case "isSyncing":
				var responseBody struct {
					Data NodeSyncingStatusDto `json:"data"`
				}
				err = json.Unmarshal(resp.data, &responseBody)
				if err != nil {
					c.WriteJSON(fiber.Map{
						"error": err.Error(),
					})
					break
				}
				newNodeResponse.Syncing.HeadSlot = responseBody.Data.HeadSlot
				newNodeResponse.Syncing.SyncDistance = responseBody.Data.SyncDistance
				newNodeResponse.Syncing.IsSyncing = responseBody.Data.IsSyncing
				newNodeResponse.Syncing.IsOptimistic = responseBody.Data.IsOptimistic
				break
			}
		}
		close(results)

		c.WriteJSON(newNodeResponse)

		time.Sleep(time.Second * 3)
	}
}

// worker is a  collection of threads for the beacon node stats
func worker(jobs <-chan request, results chan<- result) {
	chanRes := result{}
	for job := range jobs {
		chanRes.name = job.name

		client := http.Client{
			Timeout: 4 * time.Second,
		}
		req, err := http.NewRequest(http.MethodGet, job.url, bytes.NewReader([]byte(nil)))
		if err != nil {
			chanRes.err = err
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			chanRes.err = err
			return
		}

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			chanRes.err = err
			return
		}
		chanRes.data = responseData
		results <- chanRes
	}
}

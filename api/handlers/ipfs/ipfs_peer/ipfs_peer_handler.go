// Package ipfs_peer handler is the representation layer for the  ipfs peer
// it communicate  the ipfs_peer_service for business operations for the ipfs peer
package ipfs_peer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/internal/ipfs/ipfs_peer"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	nameKeyword = "name"
)

var (
	service   = ipfs_peer.NewIpfsPeerService()
	k8sClient = k8s.NewClientService()
)

// Get gets a single IPFS peer by name
// 1-get the node validated from ValidatePeerExist method
// 2-marshall node to dto and format the response
func Get(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.Peer)

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(ipfs_peer.PeerDto).FromIPFSPeer(peer)))
}

// List returns all IPFS peers
// 1-get the pagination qs default to 0
// 2-call service to return peers list
// 3-make the pagination
// 3-marshall peers to the dto struct and format the response using NewResponse
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	peers, err := service.List(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(peers.Items)))

	start, end := shared.Page(uint(len(peers.Items)), uint(page), uint(limit))
	sort.Slice(peers.Items[:], func(i, j int) bool {
		return peers.Items[j].CreationTimestamp.Before(&peers.Items[i].CreationTimestamp)
	})

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(ipfs_peer.PeerListDto).FromIPFSPeer(peers.Items[start:end])))
}

// Create creates IPFS peer from spec
// 1-Todo validate request body and return validation error
// 2-call  ipfs peer  service to create ipfs peer
// 2-marshall node to dto and format the response
func Create(c *fiber.Ctx) error {
	dto := new(ipfs_peer.PeerDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	dto.Namespace = c.Locals("namespace").(string)

	err := dto.MetaDataDto.Validate()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	peer, err := service.Create(dto)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(ipfs_peer.PeerDto).FromIPFSPeer(peer)))
}

// Delete deletes IPFS peer by name
// 1-get node from locals which checked and assigned by ValidatePeerExist
// 2-call service to delete the node
// 3-return ok if deleted with no errors
func Delete(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.Peer)

	err := service.Delete(peer)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates IPFS peer by name from spec
// 1-todo validate request body and return validation errors if exits
// 2-get node from locals which checked and assigned by ValidatePeerExist
// 3-call ipfs peer  service to update node which returns *ipfsv1alpha1.Peer
// 4-marshall node to node dto and format the response
func Update(c *fiber.Ctx) error {
	dto := new(ipfs_peer.PeerDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	peer := c.Locals("peer").(*ipfsv1alpha1.Peer)

	peer, err := service.Update(dto, peer)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(ipfs_peer.PeerDto).FromIPFSPeer(peer)))
}

// Count returns total number of peers
// 1-call  service to get length of exiting peers items
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

// ValidatePeerExist validate peer by name exist acts as a validation for all handlers the needs to find ipfs peer by name
// 1-call service to check if node exits
// 2-return 404 if it's not
// 3-save the peer to local with the key peer to be used by the other handlers
func ValidatePeerExist(c *fiber.Ctx) error {
	nameSpacedName := types.NamespacedName{
		Name:      c.Params(nameKeyword),
		Namespace: c.Locals("namespace").(string),
	}

	peer, err := service.Get(nameSpacedName)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("peer", peer)

	return c.Next()
}

func Stats(c *websocket.Conn) {
	defer c.Close()

	name := c.Params("name")
	peer := &ipfsv1alpha1.Peer{}
	nameSpacedName := types.NamespacedName{
		Namespace: c.Locals("namespace").(string),
		Name:      name,
	}

	for {

		err := k8sClient.Get(context.Background(), nameSpacedName, peer)
		if errors.IsNotFound(err) {
			c.WriteJSON(fiber.Map{
				"error": fmt.Sprintf("peer by name %s doesn't exist", name),
			})
			return
		}

		if !peer.Spec.API {
			c.WriteJSON(fiber.Map{
				"error": "peer api is not enabled",
			})
			time.Sleep(time.Second * 3)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%d/api/v0/swarm/peers", peer.Spec.APIHost, peer.Spec.APIPort), bytes.NewReader([]byte(nil)))
		if err != nil {
			c.WriteJSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}
		client := http.Client{
			Timeout: 30 * time.Second,
		}

		res, err := client.Do(req)
		if err != nil {
			c.WriteJSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}

		responseData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.WriteJSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}

		var responseBody map[string][]interface{}
		intErr := json.Unmarshal(responseData, &responseBody)
		if intErr != nil {
			c.WriteJSON(fiber.Map{
				"error": err.Error(),
			})
			return
		}

		fmt.Println(len(responseBody["Peers"]))
		time.Sleep(time.Second * 3)
	}
}

// Package ipfs_cluster_peer handler is the representation layer for the  ipfs peer
// it communicate  the ipfs_cluster_peer_service for business operations for the ipfs cluster peer
package ipfs_cluster_peer

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/community-api/internal/ipfs/ipfs_cluster_peer"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/shared"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
)

const (
	nameKeyword = "name"
)

var service = ipfs_cluster_peer.NewIpfsClusterPeerService()

// Get gets a single IPFS cluster peer by name
// 1-get the node validated from ValidateClusterPeerExist method
// 2-marshall node to dto and format the response
func Get(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.ClusterPeer)

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(ipfs_cluster_peer.ClusterPeerDto).FromIPFSClusterPeer(*peer)))
}

// List returns all IPFS cluster peers
// 1-get the pagination qs default to 0
// 2-call service to return peers list
// 3-make the pagination
// 3-marshall cluster peers to the dto struct and format the response using NewResponse
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	peers, err := service.List(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(peers.Items)))

	start, end := shared.Page(uint(len(peers.Items)), uint(page), uint(limit))
	sort.Slice(peers.Items[:], func(i, j int) bool {
		return peers.Items[j].CreationTimestamp.Before(&peers.Items[i].CreationTimestamp)
	})

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(ipfs_cluster_peer.ClusterPeerListDto).FromIPFSClusterPeer(peers.Items[start:end])))
}

// Create creates IPFS cluster peer from spec
// 1-Todo validate request body and return validation error
// 2-call ipfs cluster peer  service to create ipfs peer
// 2-marshall node to dto and format the response
func Create(c *fiber.Ctx) error {
	dto := new(ipfs_cluster_peer.ClusterPeerDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.StatusCode()).JSON(badReq)
	}

	dto.Namespace = c.Locals("namespace").(string)

	err := dto.MetaDataDto.Validate()
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	peer, err := service.Create(*dto)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(ipfs_cluster_peer.ClusterPeerDto).FromIPFSClusterPeer(peer)))
}

// Delete deletes IPFS cluster peer by name
// 1-get node from locals which checked and assigned by ValidateClusterPeerExist
// 2-call service to delete the node
// 3-return ok if deleted with no errors
func Delete(c *fiber.Ctx) error {
	peer := c.Locals("peer").(ipfsv1alpha1.ClusterPeer)

	err := service.Delete(&peer)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates IPFS cluster peer by name from spec
// 1-todo validate request body and return validation errors if exits
// 2-get node from locals which checked and assigned by ValidateClusterPeerExist
// 3-call ipfs cluster peer  service to update node which returns *ipfsv1alpha1.ClusterPeer
// 4-marshall node to node dto and format the response
func Update(c *fiber.Ctx) error {
	dto := new(ipfs_cluster_peer.ClusterPeerDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.StatusCode()).JSON(badReq)
	}

	peer := c.Locals("peer").(ipfsv1alpha1.ClusterPeer)

	err := service.Update(*dto, &peer)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(ipfs_cluster_peer.ClusterPeerDto).FromIPFSClusterPeer(peer)))
}

// Count returns total number of cluster peers
// 1-call  service to get length of exiting cluster peers items
// 2-create X-Total-Count header with the length
// 3-return
func Count(c *fiber.Ctx) error {
	length, err := service.Count(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", length))

	return c.SendStatus(http.StatusOK)
}

// ValidateClusterPeerExist validate cluster peer by name exist
// 1-call service to check if node exits
// 2-return 404 if it's not
// 3-save the peer to local with the key peer to be used by the other handlers
func ValidateClusterPeerExist(c *fiber.Ctx) error {
	nameSpacedName := types.NamespacedName{
		Name:      c.Params(nameKeyword),
		Namespace: c.Locals("namespace").(string),
	}

	peer, err := service.Get(nameSpacedName)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	c.Locals("peer", peer)

	return c.Next()
}

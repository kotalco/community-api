package beacon_node

import (
	"fmt"
	"github.com/kotalco/api/internal/ethereum2/beacon_node"
	restError "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/shared"
	"net/http"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
)

var service = beacon_node.BeaconNodeService

// Get gets a single ethereum 2.0 beacon node by name
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereum2v1alpha1.BeaconNode)

	return c.JSON(shared.NewResponse(new(beacon_node.BeaconNodeDto).FromEthereum2BeaconNode(node)))
}

// List returns all ethereum 2.0 beacon nodes
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

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(beacon_node.BeaconNodeListDto).FromEthereum2BeaconNode(nodes.Items[start:end])))
}

// Create creates ethereum 2.0 beacon node from spec
func Create(c *fiber.Ctx) error {
	dto := new(beacon_node.BeaconNodeDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restError.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(err)
	}

	node, err := service.Create(dto)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(beacon_node.BeaconNodeDto).FromEthereum2BeaconNode(node)))
}

// Delete deletes ethereum 2.0 beacon node by name
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*ethereum2v1alpha1.BeaconNode)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates ethereum 2.0 beacon node by name from spec
func Update(c *fiber.Ctx) error {
	dto := new(beacon_node.BeaconNodeDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restError.NewBadRequestError("invalid reqeust body")
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
func Count(c *fiber.Ctx) error {
	length, err := service.Count()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", *length))

	return c.SendStatus(http.StatusOK)
}

// ValidateBeaconNodeExist validate node by name exist
func ValidateBeaconNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")

	node, err := service.Get(name)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("node", node)
	return c.Next()
}

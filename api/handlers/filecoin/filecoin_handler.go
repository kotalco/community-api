package filecoin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/internal/filecoin"
	restErrors "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/shared"
	filecoinv1alpha1 "github.com/kotalco/kotal/apis/filecoin/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
)

const (
	nameKeyword      = "name"
	namespaceKeyword = "namespace"
	defaultNamespace = "default"
)

var service = filecoin.FilecoinService

// Get gets a single Filecoin node by name
// 1-get the node validated from ValidateNodeExist method
// 2-marshall node to dto and format the response
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*filecoinv1alpha1.Node)

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(filecoin.FilecoinDto).FromFilecoinNode(node)))
}

// List returns all Filecoin nodes
// 1-get the pagination qs default to 0
// 2-call service to return node models
// 3-make the pagination
// 3-marshall nodes  to Filecoin dto and format the response using NewResponse
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))

	nodes, err := service.List(c.Query(namespaceKeyword, defaultNamespace))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	start, end := shared.Page(uint(len(nodes.Items)), uint(page))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(filecoin.FilecoinListDto).FromFilecoinNode(nodes.Items[start:end])))
}

// Create creates Filecoin node from spec
// 1-Todo validate request body and return validation error
// 2-call filecoin service to create filecoin node
// 2-marshall node to dto and format the response
func Create(c *fiber.Ctx) error {
	dto := new(filecoin.FilecoinDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	node, err := service.Create(dto)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(filecoin.FilecoinDto).FromFilecoinNode(node)))
}

// Delete deletes Filecoin node by name
// 1-get node from locals which checked and assigned by ValidateNodeExist
// 2-call filecoin service to delete the node
// 3-return ok if deleted with no errors
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*filecoinv1alpha1.Node)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates Filecoin node by name from spec
// 1-todo validate request body and return validation errors if exits
// 2-get node from locals which checked and assigned by ValidateNodeExist
// 3-call filecoin service to update node which returns *filecoinv1alpha1.Node
// 4-marshall node to node dto and format the response
func Update(c *fiber.Ctx) error {
	dto := new(filecoin.FilecoinDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	node := c.Locals("node").(*filecoinv1alpha1.Node)

	node, err := service.Update(dto, node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(filecoin.FilecoinDto).FromFilecoinNode(node)))
}

// Count returns total number of nodes
// 1-call filecoin service to get exiting node list
// 2-create X-Total-Count header with the length
// 3-return
func Count(c *fiber.Ctx) error {
	length, err := service.Count(c.Query(namespaceKeyword, defaultNamespace))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", *length))

	return c.SendStatus(http.StatusOK)
}

// ValidateNodeExist  validate node by name exist acts as a validation for all handlers the needs to find filecoin node by name
// 1-call filecoin service to check if node exits
// 2-return Not found if it's not
// 3-save the node to local with the key node to be used by the other handlers
func ValidateNodeExist(c *fiber.Ctx) error {
	nameSpacedName := types.NamespacedName{
		Name:      c.Params(nameKeyword),
		Namespace: c.Query(namespaceKeyword, defaultNamespace),
	}

	node, err := service.Get(nameSpacedName)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("node", node)

	return c.Next()
}

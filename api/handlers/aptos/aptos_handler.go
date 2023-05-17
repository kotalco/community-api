package aptos

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/community-api/internal/aptos"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/shared"
	aptosv1alpha1 "github.com/kotalco/kotal/apis/aptos/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
)

const (
	nameKeyword = "name"
)

var service = aptos.NewAptosService()

// Get returns a single aptos node by name
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(aptosv1alpha1.Node)
	return c.JSON(shared.NewResponse(new(aptos.AptosDto).FromAptosNode(node)))
}

// List returns all aptos nodes
func List(c *fiber.Ctx) error {
	// default page to 0
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	nodeList, err := service.List(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	start, end := shared.Page(uint(len(nodeList.Items)), uint(page), uint(limit))
	sort.Slice(nodeList.Items[:], func(i, j int) bool {
		return nodeList.Items[j].CreationTimestamp.Before(&nodeList.Items[i].CreationTimestamp)
	})

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodeList.Items)))

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(aptos.AptosListDto).FromAptosNode(nodeList.Items[start:end])))
}

// Create created aptos node from given specs
func Create(c *fiber.Ctx) error {
	dto := new(aptos.AptosDto)
	if err := c.BodyParser(dto); err != nil {
		badReqErr := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReqErr.StatusCode()).JSON(badReqErr)
	}

	dto.Namespace = c.Locals("namespace").(string)
	if err := dto.MetaDataDto.Validate(); err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	node, err := service.Create(*dto)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}
	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(aptos.AptosDto).FromAptosNode(node)))
}

// Update updates a single aptos node by name from spec
func Update(c *fiber.Ctx) error {
	dto := new(aptos.AptosDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.StatusCode()).JSON(badReq)
	}

	node := c.Locals("node").(aptosv1alpha1.Node)

	err := service.Update(*dto, &node)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(aptos.AptosDto).FromAptosNode(node)))
}

// Count returns total number of nodes
func Count(c *fiber.Ctx) error {
	length, err := service.Count(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", length))

	return c.SendStatus(http.StatusOK)
}

// Delete a single aptos node by name
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(aptosv1alpha1.Node)

	err := service.Delete(&node)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

func ValidateNodeExist(c *fiber.Ctx) error {
	nameSpacedName := types.NamespacedName{
		Name:      c.Params(nameKeyword),
		Namespace: c.Locals("namespace").(string),
	}

	node, err := service.Get(nameSpacedName)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	c.Locals("node", node)
	return c.Next()
}

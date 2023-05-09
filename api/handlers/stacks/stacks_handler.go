package stacks

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/community-api/internal/stacks"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/shared"
	stacksv1alpha1 "github.com/kotalco/kotal/apis/stacks/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
)

const (
	nameKeyword = "name"
)

var (
	service = stacks.NewStacksService()
)

// Create creates stacks node from spec
func Create(c *fiber.Ctx) error {
	dto := new(stacks.StacksDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.StatusCode()).JSON(badReq)
	}

	dto.Namespace = c.Locals("namespace").(string)

	err := dto.MetaDataDto.Validate()
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	node, err := service.Create(*dto)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(stacks.StacksDto).FromStacksNode(node)))
}

// Get returns a single stacks node by name
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(stacksv1alpha1.Node)
	return c.JSON(shared.NewResponse(new(stacks.StacksDto).FromStacksNode(node)))
}

// List returns all stacks nodes
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

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(stacks.StacksListDto).FromStacksNode(nodeList.Items[start:end])))
}

// Update updates a single stacks node by name from spec
func Update(c *fiber.Ctx) error {
	dto := new(stacks.StacksDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.StatusCode()).JSON(badReq)
	}

	node := c.Locals("node").(stacksv1alpha1.Node)

	err := service.Update(*dto, &node)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(stacks.StacksDto).FromStacksNode(node)))
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

// Delete a single stacks node by name
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(stacksv1alpha1.Node)

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

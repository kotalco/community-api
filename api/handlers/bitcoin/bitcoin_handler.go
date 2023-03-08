package bitcoin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/community-api/internal/bitcoin"
	"github.com/kotalco/community-api/internal/core/secret"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	bitcointv1alpha1 "github.com/kotalco/kotal/apis/bitcoin/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
)

const (
	nameKeyword = "name"
)

var (
	service       = bitcoin.NewBitcoinService()
	secretService = secret.NewSecretService()
)

// Get returns a single bitcoin node by name
// 1-get the node validated from ValidateNodeExist method
// 2-marshall node to dto and format the response
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*bitcointv1alpha1.Node)
	return c.JSON(shared.NewResponse(new(bitcoin.BitcoinDto).FromBitcoinNode(node)))
}

// List returns all bitcoin nodes
// 1-get the pagination qs default to 0
// 2-call service to return node models
// 3-make the pagination
// 3-marshall nodes  to bitcoin dto and format the response using NewResponse
func List(c *fiber.Ctx) error {
	// default page to 0
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	nodeList, err := service.List(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	start, end := shared.Page(uint(len(nodeList.Items)), uint(page), uint(limit))
	sort.Slice(nodeList.Items[:], func(i, j int) bool {
		return nodeList.Items[j].CreationTimestamp.Before(&nodeList.Items[i].CreationTimestamp)
	})

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodeList.Items)))

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(bitcoin.BitcoinListDto).FromBitcoinNode(nodeList.Items[start:end])))
}

// Create created bitcoin node from given specs
// 1-call bitcoin service to create node
// 2-marshall node to bitcoinDto and format the response
func Create(c *fiber.Ctx) error {
	dto := new(bitcoin.BitcoinDto)
	if err := c.BodyParser(dto); err != nil {
		badReqErr := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReqErr.Status).JSON(badReqErr)
	}

	dto.Namespace = c.Locals("namespace").(string)
	if err := dto.MetaDataDto.Validate(); err != nil {
		return c.Status(err.Status).JSON(err)
	}

	//check for bitcoin json rpc default user secret
	rpcSec, err := secretService.Get(types.NamespacedName{
		Name:      bitcoin.BitcoinJsonRpcDefaultUserPasswordName,
		Namespace: "kotal",
	})
	if err != nil {
		if err.Status != http.StatusNotFound {
			return c.Status(err.Status).JSON(err)
		}
		//create bitcoin user default secret
		rpcSec, err = secretService.Create(&secret.SecretDto{
			MetaDataDto: k8s.MetaDataDto{Name: bitcoin.BitcoinJsonRpcDefaultUserPasswordName, Namespace: "kotal"},
			Type:        "password",
			Data:        map[string]string{"password": bitcoin.BitcoinJsonRpcDefaultUserPasswordSecret},
		})
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}
	}

	node, err := service.Create(dto, rpcSec)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}
	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(bitcoin.BitcoinDto).FromBitcoinNode(node)))
}

// Update updates a single bitcoin node by name from spec
// 1-get node from locals which checked and assigned by ValidateNodeExist
// 3-call bitcoin service to update node which returns *bitcointv1alpha1.Node
// 4-marshall node to node dto and format the response
func Update(c *fiber.Ctx) error {
	dto := new(bitcoin.BitcoinDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(err)
	}

	node := c.Locals("node").(*bitcointv1alpha1.Node)

	node, err := service.Update(dto, node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(bitcoin.BitcoinDto).FromBitcoinNode(node)))
}

// Count returns total number of nodes
// 1-call bitcoin service to get exiting node list
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

// Delete a single bitcoin node by name
// 1-get node from locals which checked and assigned by ValidateNodeExist
// 2-call service to delete the node
// 3-return ok if deleted with no errors
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*bitcointv1alpha1.Node)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
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
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("node", node)
	return c.Next()
}

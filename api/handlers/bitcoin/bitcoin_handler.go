package bitcoin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/community-api/internal/bitcoin"
	"github.com/kotalco/community-api/pkg/shared"
	bitcointv1alpha1 "github.com/kotalco/kotal/apis/bitcoin/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	nameKeyword = "name"
)

var service = bitcoin.NewBitcoinService()

func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*bitcointv1alpha1.Node)
	return c.JSON(shared.NewResponse(new(bitcoin.BitcoinDto).FromBitcoinNode(node)))
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

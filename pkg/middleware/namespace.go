package middleware

import (
	"github.com/gofiber/fiber/v2"
	"k8s.io/apimachinery/pkg/types"
)

func Namespace(c *fiber.Ctx) error {
	namespace:=c.Query("namespace","default")

	name := c.Params("name")

	namespacedName:=types.NamespacedName{
		Name: name,
		Namespace: namespace,
	}
	c.Locals("namespacedName", namespacedName)

	c.Next()
	return nil
}


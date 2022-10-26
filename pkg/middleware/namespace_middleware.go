package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func IsNamespace(c *fiber.Ctx) error {
	namespace := c.Locals("namespace")
	if namespace == nil {
		c.Locals("namespace", "default")
	}
	return c.Next()
}

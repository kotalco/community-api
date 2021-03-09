package handlers

import "github.com/gofiber/fiber/v2"

// Handler is Ethereum node handler interface
type Handler interface {
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Register(fiber.Router)
}

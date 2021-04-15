package handlers

import "github.com/gofiber/fiber/v2"

// Handler is request handler interface
// Handler accepts and returns models
// Handler process and store nodes
type Handler interface {
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Register(fiber.Router)
}

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// ValidatorHandler is Ethereum 2.0 validator client handler
type ValidatorHandler struct{}

// NewValidatorHandler creates a new Ethereum 2.0 validator client handler
func NewValidatorHandler() handlers.Handler {
	return &ValidatorHandler{}
}

// Get gets a single Ethereum 2.0 validator client by name
func (p *ValidatorHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a validator client")
}

// List returns all Ethereum 2.0 validator clients
func (p *ValidatorHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all validator clients")
}

// Create creates Ethereum 2.0 validator client from spec
func (p *ValidatorHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a validator client")
}

// Delete deletes Ethereum 2.0 validator client by name
func (p *ValidatorHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a validator client")
}

// Update updates Ethereum 2.0 validator client by name from spec
func (p *ValidatorHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a validator client")
}

// Register registers all handlers on the given router
func (p *ValidatorHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

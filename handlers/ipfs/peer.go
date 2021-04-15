package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

// PeerHandler is IPFS peer handler
type PeerHandler struct{}

func NewPeerHandler() handlers.Handler {
	return &PeerHandler{}
}

func (p *PeerHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a peer")
}

func (p *PeerHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all peers")
}

func (p *PeerHandler) Create(c *fiber.Ctx) error {
	return c.SendString("Create a peer")
}

func (p *PeerHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a peer")
}

func (p *PeerHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a peer")
}

func (p *PeerHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

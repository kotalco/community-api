package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
)

func main() {
	app := fiber.New()

	api := app.Group("api")
	v1 := api.Group("v1")
	ethereum := v1.Group("ethereum")
	nodes := ethereum.Group("nodes")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Kotal API")
	})

	ethereumHandler := handlers.EthereumHandler{}
	ethereumHandler.Register(nodes)

	app.Listen(":3000")
}

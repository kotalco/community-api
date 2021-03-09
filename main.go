package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	handlers "github.com/kotalco/api/handlers/ethereum"
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

	var nodeHandler handlers.Handler
	if os.Getenv("MOCK") == "true" {
		nodeHandler = handlers.NewNodeMockHandler()
	} else {
		nodeHandler = handlers.NewNodeHandler()
	}

	nodeHandler.Register(nodes)

	app.Listen(":3000")
}

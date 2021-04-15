package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kotalco/api/handlers"
	ethereumHandlers "github.com/kotalco/api/handlers/ethereum"
)

func main() {
	app := fiber.New()

	// register middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// routing groups
	api := app.Group("api")
	v1 := api.Group("v1")
	ethereum := v1.Group("ethereum")
	nodes := ethereum.Group("nodes")

	// register handlers
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Kotal API")
	})

	var nodeHandler handlers.Handler
	if os.Getenv("MOCK") == "true" {
		nodeHandler = ethereumHandlers.NewNodeMockHandler()
	} else {
		nodeHandler = ethereumHandlers.NewNodeHandler()
	}

	nodeHandler.Register(nodes)

	app.Listen(":3000")
}

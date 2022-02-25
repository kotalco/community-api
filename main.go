package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kotalco/api/api"
	"github.com/kotalco/api/pkg/configs"
	"github.com/kotalco/api/pkg/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	config := configs.FiberConfig()
	app := fiber.New(config)

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	api.MapUrl(app)

	server.StartServerWithGracefulShutdown(app)
}

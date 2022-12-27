package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kotalco/community-api/api"
	"github.com/kotalco/community-api/pkg/configs"
	logger2 "github.com/kotalco/community-api/pkg/logger"
	"github.com/kotalco/community-api/pkg/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	logger2.Warn("MAIN", errors.New("some error"))
	logger2.Info("MAIN", "message")
	logger2.Error("MAIN", errors.New("some aerror"))
	config := configs.FiberConfig()
	app := fiber.New(config)

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	api.MapUrl(app)

	server.StartServerWithGracefulShutdown(app)
}

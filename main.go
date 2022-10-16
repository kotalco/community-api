package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/kotalco/community-api/api"
	"github.com/kotalco/community-api/pkg/configs"
	"github.com/kotalco/community-api/pkg/server"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	fmt.Println(traefikv1alpha1.IngressRoute{})
	config := configs.FiberConfig()
	app := fiber.New(config)

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	api.MapUrl(app)

	server.StartServerWithGracefulShutdown(app)
}

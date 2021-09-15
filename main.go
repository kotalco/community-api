package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	coreHandlers "github.com/kotalco/api/handlers/core"
	ethereumHandlers "github.com/kotalco/api/handlers/ethereum"
	ethereum2Handlers "github.com/kotalco/api/handlers/ethereum2"
	ipfsHandlers "github.com/kotalco/api/handlers/ipfs"
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

	core := v1.Group("core")
	secrets := core.Group("secrets")
	storageClasses := core.Group("storageclasses")

	ethereum := v1.Group("ethereum")
	nodes := ethereum.Group("nodes")

	ipfs := v1.Group("ipfs")
	peers := ipfs.Group("peers")
	clusterpeers := ipfs.Group("clusterpeers")

	ethereum2 := v1.Group("ethereum2")
	beaconnodes := ethereum2.Group("beaconnodes")
	validators := ethereum2.Group("validators")

	// register handlers
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Kotal API")
	})

	nodeHandler := ethereumHandlers.NewNodeHandler()
	peerHandler := ipfsHandlers.NewPeerHandler()
	clusterPeerHandler := ipfsHandlers.NewClusterPeerHandler()
	beaconHandler := ethereum2Handlers.NewBeaconNodeHandler()
	validatorHandler := ethereum2Handlers.NewValidatorHandler()
	secretHandler := coreHandlers.NewSecretHandler()
	storageClassHandler := coreHandlers.NewStorageClassHandler()

	nodeHandler.Register(nodes)
	peerHandler.Register(peers)
	clusterPeerHandler.Register(clusterpeers)
	beaconHandler.Register(beaconnodes)
	validatorHandler.Register(validators)
	secretHandler.Register(secrets)
	storageClassHandler.Register(storageClasses)

	port := os.Getenv("KOTAL_API_SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	chainlinkHandlers "github.com/kotalco/api/handlers/chainlink"
	coreHandlers "github.com/kotalco/api/handlers/core"
	ethereumHandlers "github.com/kotalco/api/handlers/ethereum"
	ethereum2Handlers "github.com/kotalco/api/handlers/ethereum2"
	ipfsHandlers "github.com/kotalco/api/handlers/ipfs"
	polkadotHandlers "github.com/kotalco/api/handlers/polkadot"
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

	chainlink := v1.Group("chainlink")
	chainlinkNodes := chainlink.Group("nodes")

	polkadot := v1.Group("polkadot")
	polkadotNodes := polkadot.Group("nodes")

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

	chainlinkHandlers.NewNodeHandler().Register(chainlinkNodes)
	polkadotHandlers.NewNodeHandler().Register(polkadotNodes)
	ethereumHandlers.NewNodeHandler().Register(nodes)
	ipfsHandlers.NewPeerHandler().Register(peers)
	ipfsHandlers.NewClusterPeerHandler().Register(clusterpeers)
	ethereum2Handlers.NewBeaconNodeHandler().Register(beaconnodes)
	ethereum2Handlers.NewValidatorHandler().Register(validators)
	coreHandlers.NewSecretHandler().Register(secrets)
	coreHandlers.NewStorageClassHandler().Register(storageClasses)

	port := os.Getenv("KOTAL_API_SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/api/handlers/chainlink"
	"github.com/kotalco/community-api/api/handlers/core/secret"
	"github.com/kotalco/community-api/api/handlers/core/storage_class"
	"github.com/kotalco/community-api/api/handlers/ethereum"
	"github.com/kotalco/community-api/api/handlers/ethereum2/beacon_node"
	"github.com/kotalco/community-api/api/handlers/ethereum2/validator"
	"github.com/kotalco/community-api/api/handlers/filecoin"
	"github.com/kotalco/community-api/api/handlers/ipfs/ipfs_cluster_peer"
	"github.com/kotalco/community-api/api/handlers/ipfs/ipfs_peer"
	"github.com/kotalco/community-api/api/handlers/near"
	"github.com/kotalco/community-api/api/handlers/polkadot"
	"github.com/kotalco/community-api/api/handlers/shared"
	"github.com/kotalco/community-api/pkg/middleware"
)

// MapUrl abstracted function to map and register all the url for the application
// Used in the main.go
// helps to keep all the endpoints' definition in one place
// the only place to interact with handlers and middlewares
func MapUrl(app *fiber.App, handlers ...fiber.Handler) {
	// routing groups
	api := app.Group("api")
	v1 := api.Group("v1")
	for i := 0; i < len(handlers); i++ {
		v1.Use(handlers[i])
	}
	v1.Use(middleware.SetNamespace)
	// chainlink group
	chainlinkGroup := v1.Group("chainlink")
	chainlinkNodes := chainlinkGroup.Group("nodes")

	chainlinkNodes.Post("/", chainlink.Create)
	chainlinkNodes.Head("/", chainlink.Count)
	chainlinkNodes.Get("/", chainlink.List)
	chainlinkNodes.Get("/:name", chainlink.ValidateNodeExist, chainlink.Get)
	chainlinkNodes.Get("/:name/logs", websocket.New(shared.Logger))
	chainlinkNodes.Get("/:name/status", websocket.New(shared.Status))
	chainlinkNodes.Put("/:name", chainlink.ValidateNodeExist, chainlink.Update)
	chainlinkNodes.Delete("/:name", chainlink.ValidateNodeExist, chainlink.Delete)

	//ethereum group
	ethereumGroup := v1.Group("ethereum")
	ethereumNodes := ethereumGroup.Group("nodes")
	ethereumNodes.Post("/", ethereum.Create)
	ethereumNodes.Head("/", ethereum.Count)
	ethereumNodes.Get("/", ethereum.List)
	ethereumNodes.Get("/:name", ethereum.ValidateNodeExist, ethereum.Get)
	ethereumNodes.Get("/:name/logs", websocket.New(shared.Logger))
	ethereumNodes.Get("/:name/status", websocket.New(shared.Status))
	ethereumNodes.Get("/:name/stats", websocket.New(ethereum.Stats))
	ethereumNodes.Put("/:name", ethereum.ValidateNodeExist, ethereum.Update)
	ethereumNodes.Delete("/:name", ethereum.ValidateNodeExist, ethereum.Delete)

	//core group
	coreGroup := v1.Group("core")
	//secret group
	secrets := coreGroup.Group("secrets")
	secrets.Post("/", secret.Create)
	secrets.Head("/", secret.Count)
	secrets.Get("/", secret.List)
	secrets.Get("/:name", secret.ValidateSecretExist, secret.Get)
	secrets.Put("/:name", secret.ValidateSecretExist, secret.Update)
	secrets.Delete("/:name", secret.ValidateSecretExist, secret.Delete)
	//storage class group
	storageClasses := coreGroup.Group("storageclasses")
	storageClasses.Post("/", storage_class.Create)
	storageClasses.Get("/", storage_class.List)
	storageClasses.Get("/:name", storage_class.ValidateStorageClassExist, storage_class.Get)
	storageClasses.Put("/:name", storage_class.ValidateStorageClassExist, storage_class.Update)
	storageClasses.Delete("/:name", storage_class.ValidateStorageClassExist, storage_class.Delete)

	//ethereum2 group
	ethereum2 := v1.Group("ethereum2")
	//beaconnodes group
	beaconnodesGroup := ethereum2.Group("beaconnodes")
	beaconnodesGroup.Post("/", beacon_node.Create)
	beaconnodesGroup.Head("/", beacon_node.Count)
	beaconnodesGroup.Get("/", beacon_node.List)
	beaconnodesGroup.Get("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Get)
	beaconnodesGroup.Get("/:name/logs", websocket.New(shared.Logger))
	beaconnodesGroup.Get("/:name/status", websocket.New(shared.Status))
	beaconnodesGroup.Put("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Update)
	beaconnodesGroup.Delete("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Delete)
	//validators group
	validatorsGroup := ethereum2.Group("validators")
	validatorsGroup.Post("/", validator.Create)
	validatorsGroup.Head("/", validator.Count)
	validatorsGroup.Get("/", validator.List)
	validatorsGroup.Get("/:name", validator.ValidateValidatorExist, validator.Get)
	validatorsGroup.Get("/:name/logs", websocket.New(shared.Logger))
	validatorsGroup.Get("/:name/status", websocket.New(shared.Status))
	validatorsGroup.Put("/:name", validator.ValidateValidatorExist, validator.Update)
	validatorsGroup.Delete("/:name", validator.ValidateValidatorExist, validator.Delete)

	//filecoin group
	filecoinGroup := v1.Group("filecoin")
	filecoinNodes := filecoinGroup.Group("nodes")
	filecoinNodes.Post("/", filecoin.Create)
	filecoinNodes.Head("/", filecoin.Count)
	filecoinNodes.Get("/", filecoin.List)
	filecoinNodes.Get("/:name", filecoin.ValidateNodeExist, filecoin.Get)
	filecoinNodes.Get("/:name/logs", websocket.New(shared.Logger))
	filecoinNodes.Get("/:name/status", websocket.New(shared.Status))
	filecoinNodes.Put("/:name", filecoin.ValidateNodeExist, filecoin.Update)
	filecoinNodes.Delete("/:name", filecoin.ValidateNodeExist, filecoin.Delete)

	//ipfs group
	ipfsGroup := v1.Group("ipfs")
	//ipfs peer group
	ipfsPeersGroup := ipfsGroup.Group("peers")
	ipfsPeersGroup.Post("/", ipfs_peer.Create)
	ipfsPeersGroup.Head("/", ipfs_peer.Count)
	ipfsPeersGroup.Get("/", ipfs_peer.List)
	ipfsPeersGroup.Get("/:name", ipfs_peer.ValidatePeerExist, ipfs_peer.Get)
	ipfsPeersGroup.Get("/:name/logs", websocket.New(shared.Logger))
	ipfsPeersGroup.Get("/:name/status", websocket.New(shared.Status))
	ipfsPeersGroup.Put("/:name", ipfs_peer.ValidatePeerExist, ipfs_peer.Update)
	ipfsPeersGroup.Delete("/:name", ipfs_peer.ValidatePeerExist, ipfs_peer.Delete)
	//ipfs peer group
	clusterpeersGroup := ipfsGroup.Group("clusterpeers")
	clusterpeersGroup.Post("/", ipfs_cluster_peer.Create)
	clusterpeersGroup.Head("/", ipfs_cluster_peer.Count)
	clusterpeersGroup.Get("/", ipfs_cluster_peer.List)
	clusterpeersGroup.Get("/:name", ipfs_cluster_peer.ValidateClusterPeerExist, ipfs_cluster_peer.Get)
	clusterpeersGroup.Get("/:name/logs", websocket.New(shared.Logger))
	clusterpeersGroup.Get("/:name/status", websocket.New(shared.Status))
	clusterpeersGroup.Put("/:name", ipfs_cluster_peer.ValidateClusterPeerExist, ipfs_cluster_peer.Update)
	clusterpeersGroup.Delete("/:name", ipfs_cluster_peer.ValidateClusterPeerExist, ipfs_cluster_peer.Delete)

	//near group
	nearGroup := v1.Group("near")
	nearNodesGroup := nearGroup.Group("nodes")
	nearNodesGroup.Post("/", near.Create)
	nearNodesGroup.Head("/", near.Count)
	nearNodesGroup.Get("/", near.List)
	nearNodesGroup.Get("/:name", near.ValidateNodeExist, near.Get)
	nearNodesGroup.Get("/:name/logs", websocket.New(shared.Logger))
	nearNodesGroup.Get("/:name/status", websocket.New(shared.Status))
	nearNodesGroup.Get("/:name/stats", websocket.New(near.Stats))
	nearNodesGroup.Put("/:name", near.ValidateNodeExist, near.Update)
	nearNodesGroup.Delete("/:name", near.ValidateNodeExist, near.Delete)

	polkadotGroup := v1.Group("polkadot")
	polkadotNodesGroup := polkadotGroup.Group("nodes")
	polkadotNodesGroup.Post("/", polkadot.Create)
	polkadotNodesGroup.Head("/", polkadot.Count)
	polkadotNodesGroup.Get("/", polkadot.List)
	polkadotNodesGroup.Get("/:name", polkadot.ValidateNodeExist, polkadot.Get)
	polkadotNodesGroup.Get("/:name/logs", websocket.New(shared.Logger))
	polkadotNodesGroup.Get("/:name/status", websocket.New(shared.Status))
	polkadotNodesGroup.Get("/:name/stats", websocket.New(polkadot.Stats))
	polkadotNodesGroup.Put("/:name", polkadot.ValidateNodeExist, polkadot.Update)
	polkadotNodesGroup.Delete("/:name", polkadot.ValidateNodeExist, polkadot.Delete)

}

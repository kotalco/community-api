package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/api/handlers/bitcoin"
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

	chainlinkNodes.Post("/", chainlink.Create, middleware.IsDuplicated)
	chainlinkNodes.Head("/", chainlink.Count)
	chainlinkNodes.Get("/", chainlink.List)
	chainlinkNodes.Get("/:name", chainlink.ValidateNodeExist, chainlink.Get)
	chainlinkNodes.Get("/:name/logs", websocket.New(shared.Logger))
	chainlinkNodes.Get("/:name/status", websocket.New(shared.Status))
	chainlinkNodes.Get("/:name/metrics", websocket.New(shared.Metrics))
	chainlinkNodes.Put("/:name", chainlink.ValidateNodeExist, chainlink.Update)
	chainlinkNodes.Delete("/:name", chainlink.ValidateNodeExist, chainlink.Delete)

	//ethereum group
	ethereumGroup := v1.Group("ethereum")
	ethereumNodes := ethereumGroup.Group("nodes")
	ethereumNodes.Post("/", ethereum.Create, middleware.IsDuplicated)
	ethereumNodes.Head("/", ethereum.Count)
	ethereumNodes.Get("/", ethereum.List)
	ethereumNodes.Get("/:name", ethereum.ValidateNodeExist, ethereum.Get)
	ethereumNodes.Get("/:name/logs", websocket.New(shared.Logger))
	ethereumNodes.Get("/:name/status", websocket.New(shared.Status))
	ethereumNodes.Get("/:name/stats", websocket.New(ethereum.Stats))
	ethereumNodes.Get("/:name/metrics", websocket.New(shared.Metrics))
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
	storageClasses.Get("/", storage_class.List)
	storageClasses.Get("/:name", storage_class.ValidateStorageClassExist, storage_class.Get)
	storageClasses.Put("/:name", storage_class.ValidateStorageClassExist, storage_class.Update)
	storageClasses.Delete("/:name", storage_class.ValidateStorageClassExist, storage_class.Delete)

	//ethereum2 group
	ethereum2 := v1.Group("ethereum2")
	//beaconnodes group
	beaconnodesGroup := ethereum2.Group("beaconnodes")
	beaconnodesGroup.Post("/", beacon_node.Create, middleware.IsDuplicated)
	beaconnodesGroup.Head("/", beacon_node.Count)
	beaconnodesGroup.Get("/", beacon_node.List)
	beaconnodesGroup.Get("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Get)
	beaconnodesGroup.Get("/:name/logs", websocket.New(shared.Logger))
	beaconnodesGroup.Get("/:name/status", websocket.New(shared.Status))
	beaconnodesGroup.Get("/:name/metrics", websocket.New(shared.Metrics))
	beaconnodesGroup.Get("/:name/stats", websocket.New(beacon_node.Stats))
	beaconnodesGroup.Put("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Update)
	beaconnodesGroup.Delete("/:name", beacon_node.ValidateBeaconNodeExist, beacon_node.Delete)
	//validators group
	validatorsGroup := ethereum2.Group("validators")
	validatorsGroup.Post("/", validator.Create, middleware.IsDuplicated)
	validatorsGroup.Head("/", validator.Count)
	validatorsGroup.Get("/", validator.List)
	validatorsGroup.Get("/:name", validator.ValidateValidatorExist, validator.Get)
	validatorsGroup.Get("/:name/logs", websocket.New(shared.Logger))
	validatorsGroup.Get("/:name/status", websocket.New(shared.Status))
	validatorsGroup.Get("/:name/metrics", websocket.New(shared.Metrics))
	validatorsGroup.Put("/:name", validator.ValidateValidatorExist, validator.Update)
	validatorsGroup.Delete("/:name", validator.ValidateValidatorExist, validator.Delete)

	//filecoin group
	filecoinGroup := v1.Group("filecoin")
	filecoinNodes := filecoinGroup.Group("nodes")
	filecoinNodes.Post("/", filecoin.Create, middleware.IsDuplicated)
	filecoinNodes.Head("/", filecoin.Count)
	filecoinNodes.Get("/", filecoin.List)
	filecoinNodes.Get("/:name", filecoin.ValidateNodeExist, filecoin.Get)
	filecoinNodes.Get("/:name/logs", websocket.New(shared.Logger))
	filecoinNodes.Get("/:name/status", websocket.New(shared.Status))
	filecoinNodes.Get("/:name/metrics", websocket.New(shared.Metrics))
	filecoinNodes.Put("/:name", filecoin.ValidateNodeExist, filecoin.Update)
	filecoinNodes.Delete("/:name", filecoin.ValidateNodeExist, filecoin.Delete)

	//ipfs group
	ipfsGroup := v1.Group("ipfs")
	//ipfs peer group
	ipfsPeersGroup := ipfsGroup.Group("peers")
	ipfsPeersGroup.Post("/", ipfs_peer.Create, middleware.IsDuplicated)
	ipfsPeersGroup.Head("/", ipfs_peer.Count)
	ipfsPeersGroup.Get("/", ipfs_peer.List)
	ipfsPeersGroup.Get("/:name", ipfs_peer.ValidatePeerExist, ipfs_peer.Get)
	ipfsPeersGroup.Get("/:name/logs", websocket.New(shared.Logger))
	ipfsPeersGroup.Get("/:name/status", websocket.New(shared.Status))
	ipfsPeersGroup.Get("/:name/stats", websocket.New(ipfs_peer.Stats))
	ipfsPeersGroup.Get("/:name/metrics", websocket.New(shared.Metrics))
	ipfsPeersGroup.Put("/:name", ipfs_peer.ValidatePeerExist, ipfs_peer.Update)
	ipfsPeersGroup.Delete("/:name", ipfs_peer.ValidatePeerExist, ipfs_peer.Delete)
	//ipfs peer group
	clusterpeersGroup := ipfsGroup.Group("clusterpeers")
	clusterpeersGroup.Post("/", ipfs_cluster_peer.Create, middleware.IsDuplicated)
	clusterpeersGroup.Head("/", ipfs_cluster_peer.Count)
	clusterpeersGroup.Get("/", ipfs_cluster_peer.List)
	clusterpeersGroup.Get("/:name", ipfs_cluster_peer.ValidateClusterPeerExist, ipfs_cluster_peer.Get)
	clusterpeersGroup.Get("/:name/logs", websocket.New(shared.Logger))
	clusterpeersGroup.Get("/:name/status", websocket.New(shared.Status))
	clusterpeersGroup.Get("/:name/metrics", websocket.New(shared.Metrics))
	clusterpeersGroup.Put("/:name", ipfs_cluster_peer.ValidateClusterPeerExist, ipfs_cluster_peer.Update)
	clusterpeersGroup.Delete("/:name", ipfs_cluster_peer.ValidateClusterPeerExist, ipfs_cluster_peer.Delete)

	//near group
	nearGroup := v1.Group("near")
	nearNodesGroup := nearGroup.Group("nodes")
	nearNodesGroup.Post("/", near.Create, middleware.IsDuplicated)
	nearNodesGroup.Head("/", near.Count)
	nearNodesGroup.Get("/", near.List)
	nearNodesGroup.Get("/:name", near.ValidateNodeExist, near.Get)
	nearNodesGroup.Get("/:name/logs", websocket.New(shared.Logger))
	nearNodesGroup.Get("/:name/status", websocket.New(shared.Status))
	nearNodesGroup.Get("/:name/stats", websocket.New(near.Stats))
	nearNodesGroup.Get("/:name/metrics", websocket.New(shared.Metrics))
	nearNodesGroup.Put("/:name", near.ValidateNodeExist, near.Update)
	nearNodesGroup.Delete("/:name", near.ValidateNodeExist, near.Delete)

	polkadotGroup := v1.Group("polkadot")
	polkadotNodesGroup := polkadotGroup.Group("nodes")
	polkadotNodesGroup.Post("/", polkadot.Create, middleware.IsDuplicated)
	polkadotNodesGroup.Head("/", polkadot.Count)
	polkadotNodesGroup.Get("/", polkadot.List)
	polkadotNodesGroup.Get("/:name", polkadot.ValidateNodeExist, polkadot.Get)
	polkadotNodesGroup.Get("/:name/logs", websocket.New(shared.Logger))
	polkadotNodesGroup.Get("/:name/status", websocket.New(shared.Status))
	polkadotNodesGroup.Get("/:name/stats", websocket.New(polkadot.Stats))
	polkadotNodesGroup.Get("/:name/metrics", websocket.New(shared.Metrics))
	polkadotNodesGroup.Put("/:name", polkadot.ValidateNodeExist, polkadot.Update)
	polkadotNodesGroup.Delete("/:name", polkadot.ValidateNodeExist, polkadot.Delete)

	bitcoinGroup := v1.Group("bitcoin")
	bitcoinNodesGroup := bitcoinGroup.Group("nodes")
	bitcoinNodesGroup.Post("/", bitcoin.Create, middleware.IsDuplicated)
	bitcoinNodesGroup.Get("/:name", bitcoin.ValidateNodeExist, bitcoin.Get)
	bitcoinNodesGroup.Get("/", bitcoin.List)
	bitcoinNodesGroup.Head("/", bitcoin.Count)
	bitcoinNodesGroup.Put("/:name", bitcoin.ValidateNodeExist, bitcoin.Update)
	bitcoinNodesGroup.Delete("/:name", bitcoin.ValidateNodeExist, bitcoin.Delete)
	bitcoinNodesGroup.Get("/:name/logs", websocket.New(shared.Logger))
	bitcoinNodesGroup.Get("/:name/status", websocket.New(shared.Status))
	bitcoinNodesGroup.Get("/:name/metrics", websocket.New(shared.Metrics))
	bitcoinNodesGroup.Get("/:name/stats", websocket.New(bitcoin.Stats))

}

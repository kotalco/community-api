package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	models "github.com/kotalco/api/models/ipfs"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterPeerMockHandler is IPFS peer handler
type ClusterPeerMockHandler struct{}

var clusterPeersStore = map[string]*ipfsv1alpha1.ClusterPeer{}

// NewClusterPeerMockHandler creates a new IPFS cluster peer handler
func NewClusterPeerMockHandler() handlers.Handler {
	return &ClusterPeerMockHandler{}
}

// Get gets a single IPFS cluster peer by name
func (p *ClusterPeerMockHandler) Get(c *fiber.Ctx) error {
	name := c.Params("name")
	peer := clusterPeersStore[name]
	model := models.FromIPFSClusterPeer(peer)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"clusterpeer": model,
	})

}

// List returns all IPFS cluster peers
func (p *ClusterPeerMockHandler) List(c *fiber.Ctx) error {
	peers := []models.ClusterPeer{}
	for _, peer := range clusterPeersStore {
		peers = append(peers, *models.FromIPFSClusterPeer(peer))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"clusterpeers": peers,
	})
}

// Create creates IPFS cluster peer from spec
func (p *ClusterPeerMockHandler) Create(c *fiber.Ctx) error {
	model := new(models.ClusterPeer)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	if clusterPeersStore[model.Name] != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("cluster peer by name %s already exist", model.Name),
		})
	}

	peer := &ipfsv1alpha1.ClusterPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name: model.Name,
		},
	}

	if model.PeerEndpoint != "" {
		peer.Spec.PeerEndpoint = model.PeerEndpoint
	}

	if model.Consensus != "" {
		peer.Spec.Consensus = ipfsv1alpha1.ConsensusAlgorithm(model.Consensus)
	}

	if model.ID != "" {
		peer.Spec.ID = model.ID
	}

	if model.PrivatekeySecretName != "" {
		peer.Spec.PrivatekeySecretName = model.PrivatekeySecretName
	}

	if len(model.TrustedPeers) != 0 {
		peer.Spec.TrustedPeers = model.TrustedPeers
	}

	if len(model.BootstrapPeers) != 0 {
		peer.Spec.BootstrapPeers = model.BootstrapPeers
	}

	if model.ClusterSecretName != "" {
		peer.Spec.ClusterSecretName = model.ClusterSecretName
	}

	peer.Default()

	clusterPeersStore[model.Name] = peer

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"clusterpeer": models.FromIPFSClusterPeer(peer),
	})
}

// Delete deletes IPFS cluster peer by name
func (p *ClusterPeerMockHandler) Delete(c *fiber.Ctx) error {
	name := c.Params("name")
	delete(clusterPeersStore, name)
	return c.SendStatus(http.StatusNoContent)
}

// Update updates IPFS cluster peer by name from spec
func (p *ClusterPeerMockHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a mock cluster peer")
}

// validateClusterPeerExist validate cluster peer by name exist
func validateClusterPeerExist(c *fiber.Ctx) error {
	name := c.Params("name")

	if clusterPeersStore[name] != nil {
		return c.Next()
	}
	return c.Status(http.StatusNotFound).JSON(map[string]string{
		"error": fmt.Sprintf("cluster peer by name %s doesn't exist", c.Params("name")),
	})
}

// Register registers all handlers on the given router
func (p *ClusterPeerMockHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", validateClusterPeerExist, p.Get)
	router.Put("/:name", validateClusterPeerExist, p.Update)
	router.Delete("/:name", validateClusterPeerExist, p.Delete)
}

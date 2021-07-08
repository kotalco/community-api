package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/ipfs"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ClusterPeerHandler is IPFS peer handler
type ClusterPeerHandler struct{}

// NewClusterPeerHandler creates a new IPFS cluster peer handler
func NewClusterPeerHandler() handlers.Handler {
	return &ClusterPeerHandler{}
}

// Get gets a single IPFS cluster peer by name
func (p *ClusterPeerHandler) Get(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.ClusterPeer)
	model := models.FromIPFSClusterPeer(peer)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"clusterpeer": model,
	})

}

// List returns all IPFS cluster peers
func (p *ClusterPeerHandler) List(c *fiber.Ctx) error {
	peers := &ipfsv1alpha1.ClusterPeerList{}
	if err := k8s.Client().List(c.Context(), peers); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all cluster peers",
		})
	}

	peerModels := []models.ClusterPeer{}
	for _, peer := range peers.Items {
		peerModels = append(peerModels, *models.FromIPFSClusterPeer(&peer))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"clusterpeers": peerModels,
	})
}

// Create creates IPFS cluster peer from spec
func (p *ClusterPeerHandler) Create(c *fiber.Ctx) error {
	model := new(models.ClusterPeer)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	peer := &ipfsv1alpha1.ClusterPeer{
		ObjectMeta: metav1.ObjectMeta{
			Name:              model.Name,
			Namespace:         "default",
			CreationTimestamp: metav1.Now(),
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

	if model.CPU != "" {
		peer.Spec.CPU = model.CPU
	}
	if model.CPULimit != "" {
		peer.Spec.CPULimit = model.CPULimit
	}
	if model.Memory != "" {
		peer.Spec.Memory = model.Memory
	}
	if model.MemoryLimit != "" {
		peer.Spec.MemoryLimit = model.MemoryLimit
	}
	if model.Storage != "" {
		peer.Spec.Storage = model.Storage
	}

	if os.Getenv("MOCK") == "true" {
		peer.Default()
	}

	if err := k8s.Client().Create(c.Context(), peer); err != nil {
		log.Println(err)
		if errors.IsAlreadyExists(err) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": fmt.Sprintf("cluster peer by name %s already exist", model.Name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create cluster peer",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"clusterpeer": models.FromIPFSClusterPeer(peer),
	})
}

// Delete deletes IPFS cluster peer by name
func (p *ClusterPeerHandler) Delete(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.ClusterPeer)

	if err := k8s.Client().Delete(c.Context(), peer); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete cluster peer by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates IPFS cluster peer by name from spec
func (p *ClusterPeerHandler) Update(c *fiber.Ctx) error {
	name := c.Params("name")
	model := new(models.ClusterPeer)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	peer := c.Locals("peer").(*ipfsv1alpha1.ClusterPeer)

	if model.PeerEndpoint != "" {
		peer.Spec.PeerEndpoint = model.PeerEndpoint
	}

	if model.CPU != "" {
		peer.Spec.CPU = model.CPU
	}
	if model.CPULimit != "" {
		peer.Spec.CPULimit = model.CPULimit
	}
	if model.Memory != "" {
		peer.Spec.Memory = model.Memory
	}
	if model.MemoryLimit != "" {
		peer.Spec.MemoryLimit = model.MemoryLimit
	}
	if model.Storage != "" {
		peer.Spec.Storage = model.Storage
	}

	if os.Getenv("MOCK") == "true" {
		peer.Default()
	}

	if err := k8s.Client().Update(c.Context(), peer); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't update cluster peer by name %s", name),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"clusterpeer": models.FromIPFSClusterPeer(peer),
	})
}

// validateClusterPeerExist validate cluster peer by name exist
func validateClusterPeerExist(c *fiber.Ctx) error {
	name := c.Params("name")
	peer := &ipfsv1alpha1.ClusterPeer{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, peer); err != nil {
		log.Println(err)
		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("cluster peer by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get cluster peer by name %s", name),
		})
	}

	c.Locals("peer", peer)

	return c.Next()
}

// Register registers all handlers on the given router
func (p *ClusterPeerHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", validateClusterPeerExist, p.Get)
	router.Put("/:name", validateClusterPeerExist, p.Update)
	router.Delete("/:name", validateClusterPeerExist, p.Delete)
}

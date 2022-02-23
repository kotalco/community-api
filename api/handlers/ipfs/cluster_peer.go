package handlers

import (
	"fmt"
	"github.com/kotalco/api/api/handlers"
	shared2 "github.com/kotalco/api/api/handlers/shared"
	"github.com/kotalco/api/internal/models/ipfs"
	"github.com/kotalco/api/pkg/k8s"
	"github.com/kotalco/api/pkg/shared"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterPeerHandler is IPFS peer handler
type ClusterPeerHandler struct{}

// NewClusterPeerHandler creates a new IPFS cluster peer handler
func NewClusterPeerHandler() handlers.Handler {
	return &ClusterPeerHandler{}
}

// Get gets a single IPFS cluster peer by name
func (cp *ClusterPeerHandler) Get(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.ClusterPeer)
	model := models.FromIPFSClusterPeer(peer)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"clusterpeer": model,
	})

}

// List returns all IPFS cluster peers
func (cp *ClusterPeerHandler) List(c *fiber.Ctx) error {
	peers := &ipfsv1alpha1.ClusterPeerList{}
	if err := k8s.Client().List(c.Context(), peers, client.InNamespace("default")); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all cluster peers",
		})
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(peers.Items)))

	peerModels := []models.ClusterPeer{}

	page, _ := strconv.Atoi(c.Query("page"))

	start, end := shared.Page(uint(len(peers.Items)), uint(page))
	sort.Slice(peers.Items[:], func(i, j int) bool {
		return peers.Items[j].CreationTimestamp.Before(&peers.Items[i].CreationTimestamp)
	})

	for _, peer := range peers.Items[start:end] {
		peerModels = append(peerModels, *models.FromIPFSClusterPeer(&peer))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"clusterpeers": peerModels,
	})
}

// Create creates IPFS cluster peer from spec
func (cp *ClusterPeerHandler) Create(c *fiber.Ctx) error {
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
		Spec: ipfsv1alpha1.ClusterPeerSpec{
			Resources: sharedAPIs.Resources{
				StorageClass: model.StorageClass,
			},
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
		peer.Spec.PrivateKeySecretName = model.PrivatekeySecretName
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
func (cp *ClusterPeerHandler) Delete(c *fiber.Ctx) error {
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
func (cp *ClusterPeerHandler) Update(c *fiber.Ctx) error {
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

	if len(model.BootstrapPeers) != 0 {
		peer.Spec.BootstrapPeers = model.BootstrapPeers
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

// Count returns total number of cluster peers
func (pr *ClusterPeerHandler) Count(c *fiber.Ctx) error {
	peers := &ipfsv1alpha1.ClusterPeerList{}
	if err := k8s.Client().List(c.Context(), peers, client.InNamespace("default")); err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(peers.Items)))

	return c.SendStatus(http.StatusOK)
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
func (cp *ClusterPeerHandler) Register(router fiber.Router) {
	router.Post("/", cp.Create)
	router.Head("/", cp.Count)
	router.Get("/", cp.List)
	router.Get("/:name", validateClusterPeerExist, cp.Get)
	router.Get("/:name/logs", websocket.New(shared2.Logger))
	router.Get("/:name/status", websocket.New(shared2.Status))
	router.Put("/:name", validateClusterPeerExist, cp.Update)
	router.Delete("/:name", validateClusterPeerExist, cp.Delete)
}
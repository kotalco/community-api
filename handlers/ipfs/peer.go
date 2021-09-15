package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/ipfs"
	"github.com/kotalco/api/shared"
	ipfsv1alpha1 "github.com/kotalco/kotal/apis/ipfs/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PeerHandler is IPFS peer handler
type PeerHandler struct{}

// NewPeerHandler creates a new IPFS peer handler
func NewPeerHandler() handlers.Handler {
	return &PeerHandler{}
}

// Get gets a single IPFS peer by name
func (pr *PeerHandler) Get(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.Peer)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"peer": models.FromIPFSPeer(peer),
	})
}

// List returns all IPFS peers
func (pr *PeerHandler) List(c *fiber.Ctx) error {
	peers := &ipfsv1alpha1.PeerList{}
	if err := k8s.Client().List(c.Context(), peers, client.InNamespace("default")); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all peers",
		})
	}

	peerModels := []models.Peer{}

	page := c.Query("page")
	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	start, end := shared.Page(uint(len(peers.Items)), uint(p))
	sort.Slice(peers.Items[:], func(i, j int) bool {
		return peers.Items[i].CreationTimestamp.Before(&peers.Items[j].CreationTimestamp)
	})

	for _, peer := range peers.Items[start:end] {
		peerModels = append(peerModels, *models.FromIPFSPeer(&peer))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"peers": peerModels,
	})

}

// Create creates IPFS peer from spec
func (pr *PeerHandler) Create(c *fiber.Ctx) error {
	model := new(models.Peer)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	initProfiles := []ipfsv1alpha1.Profile{}
	for _, profile := range model.InitProfiles {
		initProfiles = append(initProfiles, ipfsv1alpha1.Profile(profile))
	}

	peer := &ipfsv1alpha1.Peer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: "default",
		},
		Spec: ipfsv1alpha1.PeerSpec{
			InitProfiles: initProfiles,
		},
	}

	if os.Getenv("MOCK") == "true" {
		peer.Default()
	}

	if err := k8s.Client().Create(c.Context(), peer); err != nil {
		log.Println(err)

		if errors.IsAlreadyExists(err) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": fmt.Sprintf("peer by name %s already exist", model.Name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create peer",
		})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"peer": models.FromIPFSPeer(peer),
	})
}

// Delete deletes IPFS peer by name
func (pr *PeerHandler) Delete(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.Peer)

	if err := k8s.Client().Delete(c.Context(), peer); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete peer by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates IPFS peer by name from spec
func (pr *PeerHandler) Update(c *fiber.Ctx) error {
	model := new(models.Peer)
	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	name := c.Params("name")
	peer := c.Locals("peer").(*ipfsv1alpha1.Peer)

	if model.APIPort != 0 {
		peer.Spec.APIPort = model.APIPort
	}

	if model.APIHost != "" {
		peer.Spec.APIHost = model.APIHost
	}

	if model.GatewayPort != 0 {
		peer.Spec.GatewayPort = model.GatewayPort
	}

	if model.GatewayHost != "" {
		peer.Spec.GatewayHost = model.GatewayHost
	}

	if model.Routing != "" {
		peer.Spec.Routing = ipfsv1alpha1.RoutingMechanism(model.Routing)
	}

	if len(model.Profiles) != 0 {
		profiles := []ipfsv1alpha1.Profile{}
		for _, profile := range model.Profiles {
			profiles = append(profiles, ipfsv1alpha1.Profile(profile))
		}
		peer.Spec.Profiles = profiles
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
			"error": fmt.Sprintf("can't update peer by name %s", name),
		})
	}

	updatedModel := models.FromIPFSPeer(peer)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"peer": updatedModel,
	})
}

// Count returns total number of peers
func (pr *PeerHandler) Count(c *fiber.Ctx) error {
	peers := &ipfsv1alpha1.PeerList{}
	if err := k8s.Client().List(c.Context(), peers, client.InNamespace("default")); err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("X-Total-Count", fmt.Sprintf("%d", len(peers.Items)))

	return c.SendStatus(http.StatusOK)
}

// validatePeerExist validates ipfs peer by name exist
func validatePeerExist(c *fiber.Ctx) error {
	name := c.Params("name")
	peer := &ipfsv1alpha1.Peer{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, peer); err != nil {

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("peer by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get peer by name %s", name),
		})
	}

	c.Locals("peer", peer)

	return c.Next()
}

// Register registers all handlers on the given router
func (pr *PeerHandler) Register(router fiber.Router) {
	router.Post("/", pr.Create)
	router.Head("/", pr.Count)
	router.Get("/", pr.List)
	router.Get("/:name", validatePeerExist, pr.Get)
	router.Put("/:name", validatePeerExist, pr.Update)
	router.Delete("/:name", validatePeerExist, pr.Delete)
}

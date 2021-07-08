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

// PeerHandler is IPFS peer handler
type PeerHandler struct{}

// NewPeerHandler creates a new IPFS peer handler
func NewPeerHandler() handlers.Handler {
	return &PeerHandler{}
}

// Get gets a single IPFS peer by name
func (p *PeerHandler) Get(c *fiber.Ctx) error {
	peer := c.Locals("peer").(*ipfsv1alpha1.Peer)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"peer": models.FromIPFSPeer(peer),
	})
}

// List returns all IPFS peers
func (p *PeerHandler) List(c *fiber.Ctx) error {
	peers := &ipfsv1alpha1.PeerList{}
	if err := k8s.Client().List(c.Context(), peers); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all peers",
		})
	}

	peerModels := []models.Peer{}
	for _, peer := range peers.Items {
		peerModels = append(peerModels, *models.FromIPFSPeer(&peer))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"peers": peerModels,
	})

}

// Create creates IPFS peer from spec
func (p *PeerHandler) Create(c *fiber.Ctx) error {
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
func (p *PeerHandler) Delete(c *fiber.Ctx) error {
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
func (p *PeerHandler) Update(c *fiber.Ctx) error {
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
func (p *PeerHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", validatePeerExist, p.Get)
	router.Put("/:name", validatePeerExist, p.Update)
	router.Delete("/:name", validatePeerExist, p.Delete)
}

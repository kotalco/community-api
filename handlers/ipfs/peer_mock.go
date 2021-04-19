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

// peersStore is in-memory IPFS peers store
var peersStore = map[string]*ipfsv1alpha1.Peer{}

// PeerMockHandler is IPFS mock peer handler
type PeerMockHandler struct{}

// NewMockPeerHandler creates a new mock IPFS peer handler
func NewPeerMockHandler() handlers.Handler {
	return &PeerMockHandler{}
}

// Get gets a single mock IPFS peer by name
func (p *PeerMockHandler) Get(c *fiber.Ctx) error {
	name := c.Params("name")
	peer := peersStore[name]
	model := models.FromIPFSPeer(peer)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"peer": model,
	})
}

// List returns all IPFS mock peers
func (p *PeerMockHandler) List(c *fiber.Ctx) error {
	peerModels := []models.Peer{}

	for _, peer := range peersStore {
		model := models.FromIPFSPeer(peer)
		peerModels = append(peerModels, *model)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"peers": peerModels,
	})
}

// Create creates IPFS mock peer from spec
func (p *PeerMockHandler) Create(c *fiber.Ctx) error {
	model := new(models.Peer)

	if err := c.BodyParser(model); err != nil {
		return err
	}

	if peersStore[model.Name] != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("peer by name %s already exist", model.Name),
		})
	}

	peer := &ipfsv1alpha1.Peer{
		ObjectMeta: metav1.ObjectMeta{
			Name: model.Name,
		},
	}

	peersStore[model.Name] = peer

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"peer": model,
	})
}

// Delete deletes IPFS mock peer by name
func (p *PeerMockHandler) Delete(c *fiber.Ctx) error {
	name := c.Params("name")
	delete(peersStore, name)
	return c.SendStatus(http.StatusNoContent)
}

// Update updates IPFS mock peer by name from spec
func (p *PeerMockHandler) Update(c *fiber.Ctx) error {
	model := new(models.Peer)
	if err := c.BodyParser(model); err != nil {
		return err
	}

	name := c.Params("name")
	peer := peersStore[name]

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

	updatedModel := models.FromIPFSPeer(peer)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"peer": updatedModel,
	})
}

// validatePeerExist validates ipfs peer by name exist
func validatePeerExist(c *fiber.Ctx) error {
	name := c.Params("name")

	if peersStore[name] != nil {
		return c.Next()
	}

	return c.Status(http.StatusNotFound).JSON(fiber.Map{
		"error": fmt.Sprintf("peer by name %s doesn't exist", c.Params("name")),
	})
}

// Register registers all handlers on the given router
func (p *PeerMockHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", validatePeerExist, p.Get)
	router.Put("/:name", validatePeerExist, p.Update)
	router.Delete("/:name", validatePeerExist, p.Delete)
}

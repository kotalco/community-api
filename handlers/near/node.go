package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/api/handlers"
	sharedHandlers "github.com/kotalco/api/handlers/shared"
	"github.com/kotalco/api/k8s"
	models "github.com/kotalco/api/models/near"
	"github.com/kotalco/api/shared"
	nearv1alpha1 "github.com/kotalco/kotal/apis/near/v1alpha1"
	sharedAPIs "github.com/kotalco/kotal/apis/shared"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NodeHandler is NEAR node handler
type NodeHandler struct{}

// NewNodeHandler creates a new NEAR node handler
func NewNodeHandler() handlers.Handler {
	return &NodeHandler{}
}

// Get gets a single NEAR node by name
func (n *NodeHandler) Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*nearv1alpha1.Node)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": models.FromNEARNode(node),
	})
}

// List returns all NEAR nodes
func (n *NodeHandler) List(c *fiber.Ctx) error {
	nodes := &nearv1alpha1.NodeList{}
	if err := k8s.Client().List(c.Context(), nodes, client.InNamespace("default")); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get all nodes",
		})
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	nodeModels := []models.Node{}

	page, _ := strconv.Atoi(c.Query("page"))

	start, end := shared.Page(uint(len(nodes.Items)), uint(page))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	for _, node := range nodes.Items[start:end] {
		nodeModels = append(nodeModels, *models.FromNEARNode(&node))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"nodes": nodeModels,
	})

}

// Create creates NEAR node from spec
func (n *NodeHandler) Create(c *fiber.Ctx) error {
	model := new(models.Node)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	node := &nearv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.Name,
			Namespace: "default",
		},
		Spec: nearv1alpha1.NodeSpec{
			Network: model.Network,
			Resources: sharedAPIs.Resources{
				StorageClass: model.StorageClass,
			},
		},
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Create(c.Context(), node); err != nil {
		log.Println(err)

		if errors.IsAlreadyExists(err) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": fmt.Sprintf("node by name %s already exist", model.Name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create node",
		})

	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"node": models.FromNEARNode(node),
	})
}

// Delete deletes NEAR node by name
func (n *NodeHandler) Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*nearv1alpha1.Node)

	if err := k8s.Client().Delete(c.Context(), node); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't delete node by name %s", c.Params("name")),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates NEAR node by name from spec
func (n *NodeHandler) Update(c *fiber.Ctx) error {
	model := new(models.Node)
	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	name := c.Params("name")
	node := c.Locals("node").(*nearv1alpha1.Node)

	if model.CPU != "" {
		node.Spec.CPU = model.CPU
	}
	if model.CPULimit != "" {
		node.Spec.CPULimit = model.CPULimit
	}
	if model.Memory != "" {
		node.Spec.Memory = model.Memory
	}
	if model.MemoryLimit != "" {
		node.Spec.MemoryLimit = model.MemoryLimit
	}
	if model.Storage != "" {
		node.Spec.Storage = model.Storage
	}

	if os.Getenv("MOCK") == "true" {
		node.Default()
	}

	if err := k8s.Client().Update(c.Context(), node); err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't update node by name %s", name),
		})
	}

	updatedModel := models.FromNEARNode(node)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"node": updatedModel,
	})
}

// Count returns total number of nodes
func (n *NodeHandler) Count(c *fiber.Ctx) error {
	nodes := &nearv1alpha1.NodeList{}
	if err := k8s.Client().List(c.Context(), nodes, client.InNamespace("default")); err != nil {
		log.Println(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	return c.SendStatus(http.StatusOK)
}

// validateNodeExist validates NEAR node by name exist
func validateNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")
	node := &nearv1alpha1.Node{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: "default",
	}

	if err := k8s.Client().Get(c.Context(), key, node); err != nil {

		if errors.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(map[string]string{
				"error": fmt.Sprintf("node by name %s doesn't exist", name),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("can't get node by name %s", name),
		})
	}

	c.Locals("node", node)

	return c.Next()
}

// Register registers all handlers on the given router
func (n *NodeHandler) Register(router fiber.Router) {
	router.Post("/", n.Create)
	router.Head("/", n.Count)
	router.Get("/", n.List)
	router.Get("/:name", validateNodeExist, n.Get)
	router.Get("/:name/logs", websocket.New(sharedHandlers.Logger))
	router.Get("/:name/status", websocket.New(sharedHandlers.Status))
	router.Put("/:name", validateNodeExist, n.Update)
	router.Delete("/:name", validateNodeExist, n.Delete)
}

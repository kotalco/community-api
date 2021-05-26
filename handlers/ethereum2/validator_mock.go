package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/handlers"
	models "github.com/kotalco/api/models/ethereum2"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ValidatorMockHandler is Ethereum 2.0 mock validator client handler
type ValidatorMockHandler struct{}

var validatorsStore = map[string]*ethereum2v1alpha1.Validator{}

// NewValidatorMockHandler creates a new Ethereum 2.0 mock validator client handler
func NewValidatorMockHandler() handlers.Handler {
	return &ValidatorMockHandler{}
}

// Get gets a single Ethereum 2.0 mock validator client by name
func (p *ValidatorMockHandler) Get(c *fiber.Ctx) error {
	return c.SendString("Get a mock validator client")
}

// List returns all Ethereum 2.0 mock validator clients
func (p *ValidatorMockHandler) List(c *fiber.Ctx) error {
	return c.SendString("List all mock validator clients")
}

// Create creates Ethereum 2.0 mock validator client from spec
func (p *ValidatorMockHandler) Create(c *fiber.Ctx) error {
	model := new(models.Validator)

	if err := c.BodyParser(model); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}

	validator := &ethereum2v1alpha1.Validator{
		ObjectMeta: metav1.ObjectMeta{
			Name: model.Name,
		},
	}

	validator.Default()

	validatorsStore[model.Name] = validator

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"validator": models.FromEthereum2Validator(validator),
	})

}

// Delete deletes Ethereum 2.0 mock validator client by name
func (p *ValidatorMockHandler) Delete(c *fiber.Ctx) error {
	return c.SendString("Delete a mock validator client")
}

// Update updates Ethereum 2.0 mock validator client by name from spec
func (p *ValidatorMockHandler) Update(c *fiber.Ctx) error {
	return c.SendString("Update a mock validator client")
}

// Register registers all handlers on the given router
func (p *ValidatorMockHandler) Register(router fiber.Router) {
	router.Post("/", p.Create)
	router.Get("/", p.List)
	router.Get("/:name", p.Get)
	router.Put("/:name", p.Update)
	router.Delete("/:name", p.Delete)
}

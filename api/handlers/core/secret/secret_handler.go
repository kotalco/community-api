// Package secret handler is the representation layer for the secret domain
//implements secretService for node secrets cruds
package secret

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/internal/core/secret"
	restErrors "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/shared"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sort"
	"strconv"
)

var service = secret.SecretService

// Get gets a single  secret by name
// 1-get the node validated from ValidateSecretExist method
// 2-marshall secretModel and format the reponse
func Get(c *fiber.Ctx) error {
	secretModel := c.Locals("secret").(*corev1.Secret)

	return c.Status(http.StatusOK).JSON(new(secret.SecretDto).FromCoreSecret(secretModel))
}

// List returns all k8s secrets
// 1-get the pagination and type qs
// 2-call service to return secret models
// 3-paginate the list
// 3-marshall secrets model to secrets dto and format the response using NewResponse
func List(c *fiber.Ctx) error {
	secretType := c.Query("type")
	page, _ := strconv.Atoi(c.Query("page")) // default page to 0

	secrets, err := service.List()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	start, end := shared.Page(uint(len(secrets.Items)), uint(page))
	sort.Slice(secrets.Items[:], func(i, j int) bool {
		return secrets.Items[j].CreationTimestamp.Before(&secrets.Items[i].CreationTimestamp)
	})

	var secretListDto = make([]secret.SecretDto, 0)
	for _, sec := range secrets.Items[start:end] {
		keyType := sec.Labels["kotal.io/key-type"]
		if keyType == "" || secretType != "" && keyType != secretType {
			continue
		}
		secretListDto = append(secretListDto, *secret.SecretDto{}.FromCoreSecret(&sec))
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(secretListDto))
}

// Create creates k8s secret from spec
// 1-creates dto from request
// 2-call service to create and save the secret model
// 3-marshall the model to the dto and format the response
func Create(c *fiber.Ctx) error {

	dto := new(secret.SecretDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(err)
	}

	secretModel, err := service.Create(dto)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(secret.SecretDto).FromCoreSecret(secretModel)))
}

// Delete deletes k8s secret by name
// 1-check if secrets with this name exits
// 2-call service to make the delete action
// 3-return the respective response
func Delete(c *fiber.Ctx) error {
	secretModel := c.Locals("secret").(*corev1.Secret)

	err := service.Delete(secretModel)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates k8s secret by name from spec
func Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Count returns total number of secrets
// 1-call secrets service to count secrets items
// 2-set the X-Total-Count header with default to 0
func Count(c *fiber.Ctx) error {
	length, err := service.Count()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", &length))

	return c.SendStatus(http.StatusOK)
}

// ValidateSecretExist validate secret by name exist acts as a validation for all handlers the needs to find secret by name
// 1-call secrets service to check if secret exits
// 2-return 404 if it's not
// 3-save the secret to local with the key secret to be used by the other handlers
func ValidateSecretExist(c *fiber.Ctx) error {
	name := c.Params("name")
	secretModel, err := service.Get(name)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("secret", secretModel)

	return c.Next()

}

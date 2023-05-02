package validator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/community-api/internal/ethereum2/validator"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/shared"
	ethereum2v1alpha1 "github.com/kotalco/kotal/apis/ethereum2/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
)

const (
	nameKeyword = "name"
)

var service = validator.NewValidatorService()

// Get gets a single Ethereum 2.0 validator client by name
// 1-get the node validated from ValidateNodeExist method
// 2-marshall node to dto and format the response
func Get(c *fiber.Ctx) error {
	validatorNode := c.Locals("validator").(ethereum2v1alpha1.Validator)

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(validator.ValidatorDto).FromEthereum2Validator(validatorNode)))
}

// List returns all Ethereum 2.0 validator clients
// 1-get the pagination qs default to 0
// 2-call service to return node models
// 3-make the pagination
// 3-marshall nodes  to validator dto and format the response using NewResponse
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	validatorList, err := service.List(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(validatorList.Items)))

	start, end := shared.Page(uint(len(validatorList.Items)), uint(page), uint(limit))
	sort.Slice(validatorList.Items[:], func(i, j int) bool {
		return validatorList.Items[j].CreationTimestamp.Before(&validatorList.Items[i].CreationTimestamp)
	})

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(validator.ValidatorListDto).FromEthereum2Validator(validatorList.Items[start:end])))
}

// Create creates Ethereum 2.0 validator client from spec
// 1-Todo validate request body and return validation error
// 2-call validator  service to create validator model
// 2-marshall node to dto and format the response
func Create(c *fiber.Ctx) error {
	dto := new(validator.ValidatorDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(badReq)
	}

	dto.Namespace = c.Locals("namespace").(string)

	err := dto.MetaDataDto.Validate()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	validatorNode, err := service.Create(*dto)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(validator.ValidatorDto).FromEthereum2Validator(validatorNode)))
}

// Delete deletes Ethereum 2.0 validator client by name
// 1-get node from locals which checked and assigned by ValidateNodeExist
// 2-call validator service to delete the node
// 3-return ok if deleted with no errors
func Delete(c *fiber.Ctx) error {
	validatorNode := c.Locals("validator").(ethereum2v1alpha1.Validator)

	err := service.Delete(&validatorNode)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Update updates Ethereum 2.0 validator client by name from spec
// 1-todo validate request body and return validation errors if exits
// 2-get node from locals which checked and assigned by ValidateNodeExist
// 3-call validator service to update node which returns *ethereum2v1alpha1.Validator
// 4-marshall node to node dto and format the response
func Update(c *fiber.Ctx) error {
	dto := new(validator.ValidatorDto)

	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(err)
	}

	validatorNode := c.Locals("validator").(ethereum2v1alpha1.Validator)

	err := service.Update(*dto, &validatorNode)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(validator.ValidatorDto).FromEthereum2Validator(validatorNode)))
}

// Count returns total number of validators
// 1-call validator service to get exiting node list
// 2-create X-Total-Count header with the length
// 3-return
func Count(c *fiber.Ctx) error {
	length, err := service.Count(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", length))

	return c.SendStatus(http.StatusOK)
}

// ValidateValidatorExist  validate node by name exist acts as a validation for all handlers the needs to find validator by name
// 1-call validator service to check if node exits
// 2-return 404 if it's not
// 3-save the node to local with the key node to be used by the other handlers
func ValidateValidatorExist(c *fiber.Ctx) error {
	nameSpacedName := types.NamespacedName{
		Name:      c.Params(nameKeyword),
		Namespace: c.Locals("namespace").(string),
	}

	validatorNode, err := service.Get(nameSpacedName)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("validator", validatorNode)

	return c.Next()

}

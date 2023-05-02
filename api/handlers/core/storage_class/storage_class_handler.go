package storage_class

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/community-api/internal/core/storage_class"
	"github.com/kotalco/community-api/pkg/shared"
	storagev1 "k8s.io/api/storage/v1"
	"net/http"
	"sort"
	"strconv"
)

const (
	nameKeyword = "name"
)

var service = storage_class.NewStorageClassService()

// Get gets a single k8s storage class
// 1-get the node validated from ValidateStorageClassExist method
// 2-marshall storageClass model and format the reponse
func Get(c *fiber.Ctx) error {
	storageClass := c.Locals("storage_class").(storagev1.StorageClass)

	return c.Status(http.StatusOK).JSON(new(storage_class.StorageClassDto).FromCoreStorageClass(storageClass))
}

// List returns all k8s storage classes
// 1-get the pagination qs default to 0
// 2-call service to return node models
// 3-make the pagination
// 3-marshall nodes  to storage class dto and format the response using NewResponse
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page")) // default page to 0
	limit, _ := strconv.Atoi(c.Query("limit"))

	storageClassList, err := service.List()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	start, end := shared.Page(uint(len(storageClassList.Items)), uint(page), uint(limit))
	sort.Slice(storageClassList.Items[:], func(i, j int) bool {
		return storageClassList.Items[j].CreationTimestamp.Before(&storageClassList.Items[i].CreationTimestamp)
	})

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(storageClassList.Items)))

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(storage_class.StorageClassListDto).FromCoreSecret(storageClassList.Items[start:end])))
}

// Create creates k8s storage class from spec
// todo
func Create(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Delete deletes k8s storage class by name
// todo
func Delete(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Update updates k8s storage class by name from spec
// todo
func Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// ValidateStorageClassExist validate storage class by name exist acts as a validation for all handlers the needs to find storage class by name
// 1-call storage class service to check if storage class exits
// 2-return not found if it's not
// 3-save the storage class to local with the key storage_class to be used by the other handlers
func ValidateStorageClassExist(c *fiber.Ctx) error {
	storageClass, err := service.Get(c.Params(nameKeyword))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("storage_class", storageClass)

	return c.Next()
}

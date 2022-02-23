package storage_class

import (
	"fmt"
	"github.com/kotalco/api/internal/core/storage_class"
	"github.com/kotalco/api/pkg/shared"
	"net/http"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	storagev1 "k8s.io/api/storage/v1"
)

var service = storage_class.StorageClassService

// Get gets a single k8s storage class
func Get(c *fiber.Ctx) error {
	storageClass := c.Locals("storage_class").(*storagev1.StorageClass)

	return c.Status(http.StatusOK).JSON(new(storage_class.StorageClassDto).FromCoreStorageClass(storageClass))
}

// List returns all k8s storage classes
func List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page")) // default page to 0

	storageClassList, err := service.List()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	start, end := shared.Page(uint(len(storageClassList.Items)), uint(page))
	sort.Slice(storageClassList.Items[:], func(i, j int) bool {
		return storageClassList.Items[j].CreationTimestamp.Before(&storageClassList.Items[i].CreationTimestamp)
	})

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(storageClassList.Items)))

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(storage_class.StorageClassListDto).FromCoreSecret(storageClassList.Items[start:end])))
}

// Create creates k8s storage class from spec
func Create(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Delete deletes k8s storage class by name
func Delete(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// Update updates k8s storage class by name from spec
func Update(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusMethodNotAllowed)
}

// ValidateStorageClassExist validate storage class by name exist acts as a validation for all handlers the needs to find storage class by name
// 1-call storage class service to check if storage class exits
// 2-return not found if it's not
// 3-save the storage class to local with the key storage_class to be used by the other handlers
func ValidateStorageClassExist(c *fiber.Ctx) error {
	name := c.Params("name")
	storageClass, err := service.Get(name)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("storage_class", storageClass)

	return c.Next()
}

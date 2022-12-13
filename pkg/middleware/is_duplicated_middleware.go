package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	restError "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s/statefulset"
)

var statefulService = statefulset.NewService()

func IsDuplicated(c *fiber.Ctx) error {
	if c.Request().Header.IsPost() { // check on posts verbs only while trying to create new resource
		var bodyFields map[string]interface{}
		_ = c.BodyParser(&bodyFields)
		if bodyFields["name"] != nil {
			name := bodyFields["name"].(string)
			exists, err := statefulService.Exists(name)
			if err != nil {
				return err
			}
			if exists {
				conflictErr := restError.NewConflictError(fmt.Sprintf("resource %s already exists!", name))
				return c.Status(conflictErr.Status).JSON(conflictErr)
			}
		}
	}
	return c.Next()
}

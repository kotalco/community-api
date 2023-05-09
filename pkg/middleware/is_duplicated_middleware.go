package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	restError "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s/statefulset"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
)

var statefulService = statefulset.NewService()

func IsDuplicated(c *fiber.Ctx) error {
	var bodyFields map[string]interface{}
	_ = c.BodyParser(&bodyFields)
	if bodyFields["name"] != nil {
		name := bodyFields["name"].(string)
		record, err := statefulService.Get(types.NamespacedName{
			Namespace: c.Locals("namespace", "default").(string),
			Name:      name,
		})

		if record != nil { //check if record already exist , return conflict if true
			conflictErr := restError.NewConflictError(fmt.Sprintf("resource %s already exists!", name))
			return c.Status(conflictErr.StatusCode()).JSON(conflictErr)
		}

		if err != nil { //check if err code is notFound , pass if true, throw if any other erros
			if err.StatusCode() == http.StatusNotFound {
				return c.Next()
			} else {
				return c.Status(err.StatusCode()).JSON(err)
			}
		}
	}
	return c.Next()
}

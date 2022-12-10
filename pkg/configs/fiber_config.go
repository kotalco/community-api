package configs

import (
	"github.com/gofiber/fiber/v2"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/logger"
	"strconv"
	"time"
)

func FiberConfig() fiber.Config {
	readTimeoutSecondsCount, _ := strconv.Atoi(Environment.ServerReadTimeout)
	return fiber.Config{
		ReadTimeout:  time.Second * time.Duration(readTimeoutSecondsCount),
		ErrorHandler: defaultErrorHandler,
	}
}

// defaultErrorHandler used to catch all unhandled  run time errors mainly panics
// logs errors using logger pkg
// return custom error struct using restError pkg
var defaultErrorHandler = func(c *fiber.Ctx, err error) error {
	go logger.Panic("PANICKING", err)

	internalErr := restErrors.NewInternalServerError("some thing went wrong...")

	return c.Status(internalErr.Status).JSON(internalErr)
}

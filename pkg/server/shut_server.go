package server

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func ShutServer(a *fiber.App) {
	if err := a.Shutdown(); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("Oops... Server is not shutting down! Reason: %v", err)
	}
}

package proxy

import (
	"github.com/gofiber/fiber/v2"
)

// APIs setup
func Setup(app *fiber.App) {

	group := app.Group("/api/unregister")
	group.Get("/", unregister)
}

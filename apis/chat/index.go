package chat

import (
	"github.com/gofiber/fiber/v2"
)

// APIs setup
func Setup(app *fiber.App) {
	group := app.Group("/api/chat")

	group.Post("/msg", chatMsg)
	group.Post("/filetransfer", fileReceiver)
}

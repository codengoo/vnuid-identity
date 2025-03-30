package routes

import (
	"vnuid-identity/controllers"

	"github.com/gofiber/fiber/v2"
)

func init() {}

func SetupRoutes(app *fiber.App) {
	var api = app.Group("/api")
	api.Post("/create", controllers.AddUser)
}

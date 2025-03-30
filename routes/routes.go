package routes

import (
	authController "vnuid-identity/controllers/auth"
	userController "vnuid-identity/controllers/user"

	"github.com/gofiber/fiber/v2"
)

func init() {}

func SetupRoutes(app *fiber.App) {
	var userCtrl = app.Group("/user")
	var authCtrl = app.Group("/auth")

	userCtrl.Post("/add", userController.AddUser)
	userCtrl.Post("/add_many", userController.AddMultipleUsers)
	userCtrl.Delete("/remove_many", userController.RemoveMultipleUsers)

	authCtrl.Post("/login", authController.Login)
}

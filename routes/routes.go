package routes

import (
	authController "vnuid-identity/controllers/auth"
	userController "vnuid-identity/controllers/user"
	"vnuid-identity/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func init() {}

func SetupRoutes(app *fiber.App) {
	var userCtrl = app.Group("/user")
	var authCtrl = app.Group("/auth")

	userCtrl.Post("/add", middlewares.AuthCheck("admin"), userController.AddUser)
	userCtrl.Post("/add_many", middlewares.AuthCheck("admin"), userController.AddMultipleUsers)
	userCtrl.Delete("/remove_many", middlewares.AuthCheck("admin"), userController.RemoveMultipleUsers)

	authCtrl.Post("/login", authController.Login)
	authCtrl.Get("/get_qr", authController.GetQR)
	authCtrl.Post("/verify_qr", middlewares.AuthCheck("user"), authController.VerifyQR)
	authCtrl.Get("/ws/:session", websocket.New(authController.AuthSocket))
}

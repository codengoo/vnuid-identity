package routes

import (
	authController "vnuid-identity/controllers/auth"
	nfcController "vnuid-identity/controllers/nfc"
	userController "vnuid-identity/controllers/user"
	"vnuid-identity/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func init() {}

func SetupRoutes(app *fiber.App) {
	var userCtrl = app.Group("/manage")
	var authCtrl = app.Group("/auth")
	var nfcCtrl = app.Group("/nfc")

	var ADMIN = []string{"admin"}
	var AUTH = []string{"user", "admin", "teacher"}

	nfcCtrl.Post("/add", middlewares.AuthCheck(ADMIN), nfcController.AddNFC)
	nfcCtrl.Put("/activate/:id", middlewares.AuthCheck(ADMIN), nfcController.ActivateNFC)

	userCtrl.Post("/add", middlewares.AuthCheck(ADMIN), userController.AddUser)
	userCtrl.Post("/add_many", middlewares.AuthCheck(ADMIN), userController.AddMultipleUsers)
	userCtrl.Delete("/remove_many", middlewares.AuthCheck(ADMIN), userController.RemoveMultipleUsers)

	authCtrl.Post("/login_pass", authController.Login)
	authCtrl.Post("/login_google", authController.LoginByGoogle)
	authCtrl.Post("/login_pass_2fa", authController.LoginByPass2Fa)
	authCtrl.Post("/login_auth_2fa", authController.LoginByAuth2Fa)
	authCtrl.Post("/login_nfc_2fa", authController.LoginByNFC2Fa)

	authCtrl.Post("set_authenticator", middlewares.AuthCheck(AUTH), authController.SetAuthenticator)

	authCtrl.Get("/get_qr", authController.GetQR)
	authCtrl.Post("/verify_qr", middlewares.AuthCheck(AUTH), authController.VerifyQR)
	authCtrl.Get("/ws/:session", websocket.New(authController.AuthSocket))
}

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
	var AUTH = []string{"student", "admin", "teacher"}

	nfcCtrl.Post("/add", middlewares.AuthCheck(ADMIN), nfcController.AddNFC)
	nfcCtrl.Put("/activate/:id", middlewares.AuthCheck(ADMIN), nfcController.ActivateNFC)

	userCtrl.Post("/add", middlewares.AuthCheck(ADMIN), userController.AddUser)
	userCtrl.Delete("/remove_many", middlewares.AuthCheck(ADMIN), userController.RemoveMultipleUsers)

	authCtrl.Post("/login_google", authController.LoginByGoogle)
	authCtrl.Post("/login_pass", authController.LoginByPass)
	authCtrl.Post("/login_nfc", authController.LoginByNFC)
	authCtrl.Post("/login_bio", authController.LoginByBio)

	authCtrl.Post("/login_pass_2fa", authController.LoginByPass2Fa)
	authCtrl.Post("/login_otp_2fa", authController.LoginByOtp2Fa)
	authCtrl.Post("/login_nfc_2fa", authController.LoginByNFC2Fa)
	authCtrl.Post("/login_qr_2fa", authController.LoginByQr2Fa)
	authCtrl.Post("/login_code_2fa", middlewares.AuthCheck(AUTH), authController.LoginByCode2Fa)
	authCtrl.Post("/login_code_2fa_accept", middlewares.AuthCheck(AUTH), authController.LoginByCode2FaAccept)

	authCtrl.Post("/login_qr_accept", middlewares.AuthCheck(AUTH), authController.LoginByQrAccept)
	authCtrl.Post("/login_qr_info", middlewares.AuthCheck(AUTH), authController.LoginByQrInfo)
	authCtrl.Post("/set_authenticator", middlewares.AuthCheck(AUTH), authController.SetAuthenticator)
	authCtrl.Post("/set_biometric", middlewares.AuthCheck(AUTH), authController.SetBiometric)
	authCtrl.Post("/check_password", middlewares.AuthCheck(AUTH), authController.CheckPassword)

	app.Use("/ws", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/login/:session", websocket.New(authController.ListenLogin))
}

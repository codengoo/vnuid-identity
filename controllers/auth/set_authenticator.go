package controllers

import (
	"vnuid-identity/middlewares"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type SetAuthenticatorRequest struct {
	Password string `json:"password" validate:"required"`
}

func SetAuthenticator(ctx *fiber.Ctx) error {
	var data SetAuthenticatorRequest
	claims := ctx.Locals("user").(*middlewares.TokenClaim)

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON((fiber.Map{"error": err.Error(), "msg": msg}))
	}

	isValid, _ := models.VerifyPassword(claims.UID, data.Password)
	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Vnuid",
		AccountName: claims.UID,
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	png, err := qrcode.Encode(key.URL(), qrcode.Low, 256)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating QR code"})
	}

	if err := models.SetAuthenticator(claims.UID, key.Secret()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).Send(png)
}

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
	claims := ctx.Locals("user").(*middlewares.TokenClaim)
	var data SetAuthenticatorRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	isValid, _ := models.VerifyPassword(claims.UID, data.Password)
	if !isValid {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "Invalid password")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Vnuid",
		AccountName: claims.UID,
	})

	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	png, err := qrcode.Encode(key.URL(), qrcode.Low, 256)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	if err := models.SetAuthenticator(claims.UID, key.Secret()); err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).Send(png)
}

package controllers

import (
	"vnuid-identity/middlewares"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginQr2FaAcceptRequest struct {
	Token      string `json:"token" validate:"required"`
	Session    string `json:"session" validate:"required"`
	SaveDevice bool   `json:"save_device" validate:"required"`
}

var LOGIN_ACCEPT_KEY = "qr:login:accept:"

func LoginByQr2FaAccept(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)
	var data LoginQr2FaAcceptRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	// Extract thong tin tu client
	qrClaims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusBadRequest, err)
	}

	if qrClaims.UID != userClaims.UID {
		return utils.ReturnErrorMsg(ctx, fiber.StatusBadRequest, "Invalid token params")
	}

	// Set accept token
	if err := models.SetLoginAcceptSession(
		data.Session,
		models.Login2FaAcceptConfig{
			SaveDevice: data.SaveDevice,
			DeviceID:   qrClaims.DeviceID,
			UID:        qrClaims.UID,
			Method:     qrClaims.Method,
		},
	); err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}

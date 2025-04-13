package controllers

import (
	"vnuid-identity/middlewares"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByCode2FaAcceptRequest struct {
	SaveDevice bool `json:"save_device" validate:"required"`
	Code       int  `json:"code" validate:"required"`
}

func LoginByCode2FaAccept(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)
	var data LoginByCode2FaAcceptRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	// extract thong tin tu redis
	content, err := models.GetLoginCodeSession(userClaims.UID)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	if content.Code != data.Code {
		return utils.ReturnErrorMsg(ctx, fiber.StatusBadRequest, "Invalid code")
	}

	// Set accept token
	if err := models.SetLoginAcceptSession(
		content.Session,
		models.Login2FaAcceptConfig{
			SaveDevice: data.SaveDevice,
			DeviceID:   content.DeviceID,
			UID:        userClaims.UID,
		},
	); err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}

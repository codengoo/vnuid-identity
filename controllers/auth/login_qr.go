package controllers

import (
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByQrRequest struct {
	DeviceId   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

func LoginByQr(ctx *fiber.Ctx) error {
	var data LoginByQrRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	token, err := utils.GenerateTemporaryToken(utils.TmpTokenData{
		DeviceID:     data.DeviceId,
		DeviceName:   data.DeviceName,
		AllowMethods: []string{},
		Method:       "qr",
	})

	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	// Set session
	session, err := models.SetLoginSession("")
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"ws": session, "token": token})
}

package controllers

import (
	"vnuid-identity/entities"
	"vnuid-identity/helpers"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type BioLoginRequest struct {
	DeviceId   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
	BioCode    string `json:"bio_code" validate:"required"`
	UID        string `json:"uid" validate:"required"`
}

func LoginByBio(ctx *fiber.Ctx) error {
	var data BioLoginRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	valid, user := models.VerifyBioCode(data.UID, data.BioCode)
	if !valid {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "invalid nfc code")
	}

	// Generate token
	token, err := helpers.AddLoginSession(user, entities.Session{
		DeviceId:    data.DeviceId,
		DeviceName:  data.DeviceName,
		LoginMethod: "bio",
		SavedDevice: true,
	})
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(fiber.Map{"token": token})
}

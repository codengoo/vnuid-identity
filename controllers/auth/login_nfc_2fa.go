package controllers

import (
	"slices"
	"vnuid-identity/helpers"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByNFC2FaRequest struct {
	Token      string `json:"token" validate:"required"`
	SaveDevice bool   `json:"save_device"`
	NfcCode    string `json:"nfc_code" validate:"required"`
}

func LoginByNFC2Fa(ctx *fiber.Ctx) error {
	var data LoginByNFC2FaRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	claims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusUnauthorized, err)
	}

	// Validate
	if !slices.Contains(claims.AllowMethods, "nfc") {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "Invalid token params")
	}

	valid, user := models.VerifyNFC(claims.UID, data.NfcCode)
	if !valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid NFC code"})
	}

	// Generate token
	token, err := helpers.AddLoginSession(user, claims.DeviceID, data.SaveDevice, claims.Method)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

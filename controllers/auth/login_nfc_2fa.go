package controllers

import (
	"fmt"
	"slices"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByNFC2FaRequest struct {
	Token      string `json:"token" validate:"required"`
	DeviceId   string `json:"device_id" validate:"required"`
	SaveDevice bool   `json:"save_device"`
	NfcCode    string `json:"nfc_code" validate:"required"`
}

func LoginByNFC2Fa(ctx *fiber.Ctx) error {
	var data LoginByNFC2FaRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON((fiber.Map{"error": err.Error(), "msg": msg}))
	}

	claims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	if claims.DeviceID != data.DeviceId || !slices.Contains(claims.AllowMethods, "nfc") {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token params"})
	}

	valid, user := models.VerifyNFC(claims.UID, data.NfcCode)
	if !valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid NFC code"})
	}

	token, err := utils.GenerateToken(user, data.DeviceId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	// Create login session
	if _, err := models.CreateSession(data.DeviceId, user.ID, data.SaveDevice); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Create session: %s", err.Error()),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

package controllers

import (
	"slices"
	"vnuid-identity/entities"
	"vnuid-identity/helpers"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByAuthenticator2FaRequest struct {
	Token      string `json:"token" validate:"required"`
	Code       string `json:"code" validate:"required"`
	SaveDevice bool   `json:"save_device"`
}

func LoginByOtp2Fa(ctx *fiber.Ctx) error {
	var data LoginByAuthenticator2FaRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	claims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusUnauthorized, err)
	}

	// Validate
	if !slices.Contains(claims.AllowMethods, "otp") {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "Invalid token params")
	}

	valid, user := models.VerifyAuthenticator(claims.UID, data.Code)
	if !valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authentication code"})
	}

	// Generate token
	token, err := helpers.AddLoginSession(user, entities.Session{
		DeviceId:    claims.DeviceID,
		DeviceName:  claims.DeviceName,
		LoginMethod: claims.Method,
		SavedDevice: data.SaveDevice,
	})
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(fiber.Map{"token": token})
}

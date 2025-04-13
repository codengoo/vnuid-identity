package controllers

import (
	"fmt"
	"slices"
	"vnuid-identity/helpers"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type PassLogin2FaRequest struct {
	Token      string `json:"token" validate:"required"`
	Password   string `json:"password" validate:"required"`
	SaveDevice bool   `json:"save_device"`
}

func LoginByPass2Fa(ctx *fiber.Ctx) error {
	var data PassLogin2FaRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	claims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusUnauthorized, err)
	}

	// Validate
	fmt.Println(claims.AllowMethods)
	if !slices.Contains(claims.AllowMethods, "pass") {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "Invalid token params")
	}

	valid, user := models.VerifyPassword(claims.UID, data.Password)
	if !valid {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "Invalid password")
	}

	// Create token
	token, err := helpers.AddLoginSession(user, claims.DeviceID, data.SaveDevice, claims.Method)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(fiber.Map{"token": token})
}

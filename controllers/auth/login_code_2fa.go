package controllers

import (
	"fmt"
	"slices"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByCode2FaRequest struct {
	Token string `json:"token" validate:"required"`
}

var LOGIN_CODE_KEY = "qr:login:code:"

func LoginByCode2Fa(ctx *fiber.Ctx) error {
	var data LoginByCode2FaRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	// Validate
	claims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusBadRequest, err)
	}

	if !slices.Contains(claims.AllowMethods, "code") {
		return utils.ReturnError(ctx, fiber.StatusBadRequest, fmt.Errorf("invalid token params"))
	}

	// Set session
	session, err := models.SetLoginSession(claims.UID)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	// Create tmp storage to validate later
	code := utils.RandomInRange(10, 99)
	if err := models.SetLoginCodeSession(
		claims.UID,
		models.LoginByCode2FaConfig{
			Code:     code,
			Session:  session,
			DeviceID: claims.DeviceID,
			Method:   claims.Method,
		},
	); err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"code": code, "ws": session})
}

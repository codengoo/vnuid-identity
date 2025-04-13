package controllers

import (
	"slices"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByQr2FaRequest struct {
	Token string `json:"token" validate:"required"`
}

func LoginByQr2Fa(ctx *fiber.Ctx) error {
	var data LoginByQr2FaRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	// Validate
	claims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusBadRequest, err)
	}

	if !slices.Contains(claims.AllowMethods, "qr") {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "Invalid token params")
	}

	// Set session
	session, err := models.SetLoginSession(claims.UID)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"ws": session})
}

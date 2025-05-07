package controllers

import (
	"vnuid-identity/middlewares"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type CheckPasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

func CheckPassword(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)
	var data CheckPasswordRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	valid, _ := models.VerifyPassword(userClaims.UID, data.Password)
	if !valid {
		return utils.ReturnSuccess(ctx, map[string]any{"valid": false})
	}

	return utils.ReturnSuccess(ctx, map[string]any{"valid": true})
}

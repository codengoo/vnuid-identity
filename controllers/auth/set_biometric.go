package controllers

import (
	"vnuid-identity/middlewares"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type SetBiometricRequest struct {
	Password string `json:"password" validate:"required"`
}

func SetBiometric(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)
	var data SetAuthenticatorRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	valid, _ := models.VerifyPassword(userClaims.UID, data.Password)
	if !valid {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "Invalid password")
	}

	key, err := models.SetBiometric(userClaims.UID)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Ok", "key": key})
}

package controllers

import (
	"vnuid-identity/middlewares"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

func GetMe(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)
	profile, err := models.GetMe(userClaims.UID)
	if err != nil {
		return utils.ReturnErrorMsg(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.ReturnSuccess(ctx, profile)
}

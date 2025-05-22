package controllers

import (
	"vnuid-identity/entities"
	"vnuid-identity/middlewares"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type MeResponse struct {
	entities.Profile
	Role string `json:"role"`
}

func GetMe(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)
	profile, err := models.GetMe(userClaims.UID)
	if err != nil {
		return utils.ReturnErrorMsg(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.ReturnSuccess(ctx, MeResponse{
		Profile: profile,
		Role:    userClaims.Role})
}

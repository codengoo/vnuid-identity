package controllers

import (
	"vnuid-identity/middlewares"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByQr2FaInfoRequest struct {
	Token string `json:"token" validate:"required"`
}

func LoginByQrInfo(ctx *fiber.Ctx) error {
	var data LoginByQr2FaInfoRequest
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	qrClaims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusBadRequest, err)
	}

	// Verify: If Scan directly or step-2 must has valid uid
	if len(qrClaims.AllowMethods) != 0 && qrClaims.UID != userClaims.UID {
		return utils.ReturnErrorMsg(ctx, fiber.StatusBadRequest, "Invalid token params")
	}

	// extract session info
	result, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(result)
}

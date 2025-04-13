package controllers

import (
	"fmt"
	"vnuid-identity/middlewares"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByQr2FaInfoRequest struct {
	Token   string `json:"token" validate:"required"`
	Session string `json:"session" validate:"required"`
}

func LoginByQr2FaInfo(ctx *fiber.Ctx) error {
	var data LoginByQr2FaInfoRequest
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	qrClaims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusBadRequest, err)
	}

	// Verify
	fmt.Println(qrClaims.UID, userClaims.UID)
	if qrClaims.UID != userClaims.UID {
		return utils.ReturnErrorMsg(ctx, fiber.StatusBadRequest, "Invalid token params")
	}

	// extract session info
	result, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(result)
}

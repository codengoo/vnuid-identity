package controllers

import (
	"encoding/json"
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
	if qrClaims.UID != userClaims.UID {
		return utils.ReturnErrorMsg(ctx, fiber.StatusBadRequest, "Invalid token params")
	}

	// extract session info
	var content utils.TmpTokenData
	err = json.Unmarshal([]byte(data.Token), &content)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(content)
}

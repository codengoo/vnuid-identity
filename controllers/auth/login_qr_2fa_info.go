package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"vnuid-identity/databases"
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
		return ctx.Status(fiber.StatusBadRequest).JSON((fiber.Map{"error": err.Error(), "msg": msg}))
	}

	claims, err := utils.ParseQRToken(data.Token)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if claims.UID != userClaims.UID {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid token params"})
	}

	bgctx := context.Background()
	val, err := databases.RD.Get(bgctx, fmt.Sprintf("%s%s", LOGIN_KEY, data.Session)).Result()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var content LoginByQr2FaRequest
	err = json.Unmarshal([]byte(val), &content)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	content.Token = ""
	return ctx.Status(fiber.StatusOK).JSON(content)
}

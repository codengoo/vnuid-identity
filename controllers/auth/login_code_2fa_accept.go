package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vnuid-identity/databases"
	"vnuid-identity/middlewares"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginByCode2FaAcceptRequest struct {
	SaveDevice bool `json:"save_device" validate:"required"`
	Code       int  `json:"code" validate:"required"`
}

func LoginByCode2FaAccept(ctx *fiber.Ctx) error {
	var data LoginByCode2FaAcceptRequest
	userClaims := ctx.Locals("user").(*middlewares.TokenClaim)

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON((fiber.Map{"error": err.Error(), "msg": msg}))
	}

	bgctx := context.Background()
	content_raw, err := databases.RD.Get(bgctx, fmt.Sprintf("%s%s", LOGIN_CODE_KEY, userClaims.UID)).Result()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error getting session"})
	}

	var content LoginByCode2FaConfig
	err = json.Unmarshal([]byte(content_raw), &content)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error getting session"})
	}

	if content.Code != data.Code {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid code"})
	}

	content_new, err := json.Marshal(
		LoginByQr2FaAcceptConfig{
			SaveDevice: data.SaveDevice,
			DeviceID:   content.DeviceId,
			UID:        userClaims.SID,
		})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating session"})
	}

	err = databases.RD.Set(bgctx, fmt.Sprintf("%s%s", LOGIN_ACCEPT_KEY, content.Session), content_new, 60*time.Second).Err()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating session"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}

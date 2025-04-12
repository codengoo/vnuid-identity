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

type LoginQr2FaAcceptRequest struct {
	Token      string `json:"token" validate:"required"`
	Session    string `json:"session" validate:"required"`
	SaveDevice bool   `json:"save_device" validate:"required"`
}

type LoginByQr2FaAcceptConfig struct {
	SaveDevice bool   `json:"save_device"`
	DeviceID   string `json:"device_id"`
	UID        string `json:"uid"`
}

var LOGIN_ACCEPT_KEY = "qr:login:accept:"

func LoginByQr2FaAccept(ctx *fiber.Ctx) error {
	var data LoginQr2FaAcceptRequest
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
	content, err := json.Marshal(
		LoginByQr2FaAcceptConfig{
			SaveDevice: data.SaveDevice,
			DeviceID:   claims.DeviceID,
			UID:        claims.UID,
		})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating session"})
	}

	err = databases.RD.Set(bgctx, fmt.Sprintf("%s%s", LOGIN_ACCEPT_KEY, data.Session), content, 60*time.Second).Err()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating session"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}

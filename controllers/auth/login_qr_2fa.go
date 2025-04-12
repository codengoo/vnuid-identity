package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"
	"vnuid-identity/databases"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LoginByQr2FaRequest struct {
	Token      string `json:"token" validate:"required"`
	DeviceId   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
	Location   string `json:"location" validate:"required"`
}

var LOGIN_KEY = "qr:login:"

func LoginByQr2Fa(ctx *fiber.Ctx) error {
	var data LoginByQr2FaRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON((fiber.Map{"error": err.Error(), "msg": msg}))
	}

	claims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	if claims.DeviceID != data.DeviceId || !slices.Contains(claims.AllowMethods, "qr") {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token params"})
	}

	token, err := utils.GenerateQRToken(claims.UID, data.DeviceId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	session := uuid.New().String()
	bgctx := context.Background()
	jsonSession, err := json.Marshal(data)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating session"})
	}

	err = databases.RD.Set(bgctx, fmt.Sprintf("%s%s", LOGIN_KEY, session), jsonSession, 60*time.Second).Err()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating session"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token, "ws": session})
}

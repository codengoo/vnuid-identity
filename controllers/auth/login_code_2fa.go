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

type LoginByCode2FaRequest struct {
	Token      string `json:"token" validate:"required"`
	DeviceId   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
	Location   string `json:"location" validate:"required"`
}

type LoginByCode2FaConfig struct {
	DeviceId   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
	Location   string `json:"location" validate:"required"`
	Code       int    `json:"code" validate:"required"`
	Session    string `json:"session" validate:"required"`
}

var LOGIN_CODE_KEY = "qr:login:code:"

func LoginByCode2Fa(ctx *fiber.Ctx) error {
	var data LoginByCode2FaRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON((fiber.Map{"error": err.Error(), "msg": msg}))
	}

	claims, err := utils.ParseTemporaryToken(data.Token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	if claims.DeviceID != data.DeviceId || !slices.Contains(claims.AllowMethods, "code") {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token params"})
	}

	code := utils.RandomInRange(10, 99)
	session := uuid.New().String()
	bgctx := context.Background()

	content := LoginByCode2FaConfig{
		DeviceId:   data.DeviceId,
		DeviceName: data.DeviceName,
		Location:   data.Location,
		Code:       code,
		Session:    session,
	}
	jsonSession, err := json.Marshal(content)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating session"})
	}

	databases.RD.Set(bgctx, fmt.Sprintf("%s%s", LOGIN_CODE_KEY, claims.UID), jsonSession, 60*time.Second).Err()
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"code": code, "ws": session})
}

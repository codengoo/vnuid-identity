package controllers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/skip2/go-qrcode"
)

func genToken(device_id string, session_id string) (string, error) {
	secret_key := os.Getenv("JWT_TOKEN")

	claims := jwt.MapClaims{
		"device_id": device_id,
		"session":   session_id,
		"timestamp": time.Now(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret_key))
}

func GetQR(ctx *fiber.Ctx) error {
	device_id := ctx.Query("device_id")
	session_id := ctx.Query("session")

	token, err := genToken(device_id, session_id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	png, err := qrcode.Encode(token, qrcode.Low, 256)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating QR code"})
	}

	return ctx.Status(fiber.StatusOK).Send(png)
}

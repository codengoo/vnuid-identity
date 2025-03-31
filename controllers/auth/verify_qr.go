package controllers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type VerifyQRLoginRequest struct {
	Token    string `json:"token"`
	DeviceID string `json:"device_id"`
}

type TokenLoginClaim struct {
	jwt.RegisteredClaims
}

func parseToken(tokenString string) (string, error) {
	secret_key := os.Getenv("JWT_TOKEN")

	_, err := jwt.ParseWithClaims(tokenString, &TokenLoginClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}

		// check more here
		return []byte(secret_key), nil
	})

	return "", err

}

func VerifyQR(ctx *fiber.Ctx) error {
	var data VerifyQRLoginRequest
	err := ctx.BodyParser(&data)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// verify token
	s, err := parseToken(data.Token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	fmt.Println(s)

	return ctx.SendString("OK")
}

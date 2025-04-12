package controllers

import (
	"context"
	"fmt"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"os"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/gofiber/fiber/v2"
)

type GoogleLoginRequest struct {
	TokenId  string `json:"id_token" validate:"required"`
	DeviceId string `json:"device_id" validate:"required"`
}

func verifyGoogleIDToken(token string) (*idtoken.Payload, error) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		return nil, fmt.Errorf("invalid Google Client ID")
	}

	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, token, googleClientID)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func LoginByGoogle(ctx *fiber.Ctx) error {
	var data GoogleLoginRequest
	if err := utils.GetBodyData(ctx, &data); err != nil {
		return err
	}

	// Verify google ID
	payload, err := verifyGoogleIDToken(data.TokenId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid token: " + err.Error()})
	}

	// Extract user linked with this google account
	gid := payload.Claims["sub"].(string)
	user, err := models.GetUser(gid)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Can not find user linked with this account"})
	}

	// Check if session already logged in on device_id
	isActive := models.CheckSession(data.DeviceId, user.ID)
	if isActive {
		// Generate token
		token, err := utils.GenerateToken(user, data.DeviceId)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
		}

		// Create login session
		if _, err := models.CreateSession(data.DeviceId, user.ID, true); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Create session: %s", err.Error()),
			})
		}

		return ctx.JSON(fiber.Map{"token": token})
	} else {
		var allowList []string = []string{"password", "qr", "otp, nfc,auth"}

		token, err := utils.GenerateTemporaryToken(gid, data.DeviceId)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
		}
		return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"allow": allowList, "token": token})
	}
}

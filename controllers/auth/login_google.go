package controllers

import (
	"context"
	"fmt"
	"vnuid-identity/entities"
	"vnuid-identity/helpers"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"os"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/gofiber/fiber/v2"
)

type GoogleLoginRequest struct {
	TokenId    string `json:"id_token" validate:"required"`
	DeviceId   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
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
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}
	fmt.Println(data)

	// Verify google ID
	payload, err := verifyGoogleIDToken(data.TokenId)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusBadRequest, err)
	}

	// Extract user linked with this google account
	gid := payload.Claims["sub"].(string)
	user, err := models.GetUser(gid)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusUnauthorized, err)
	}

	// Check if session already logged in on device_id
	isActive := models.CheckSession(data.DeviceId, user.ID)
	if isActive {
		// Generate token
		token, err := helpers.AddLoginSession(user, entities.Session{
			DeviceId:    data.DeviceId,
			DeviceName:  data.DeviceName,
			LoginMethod: "google",
			SavedDevice: true,
		})
		if err != nil {
			return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
		}

		return ctx.JSON(fiber.Map{"token": token, "uid": user.ProfileId})
	} else {
		var allowList []string = []string{"pass", "qr", "code", "nfc", "otp"}

		// Gen token for step 2
		token, err := utils.GenerateTemporaryToken(
			utils.TmpTokenData{
				UID:          user.ID,
				DeviceID:     data.DeviceId,
				DeviceName:   data.DeviceName,
				AllowMethods: allowList,
				Method:       "google",
			})
		if err != nil {
			return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
		}

		return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"allow": allowList,
			"token": token,
			"uid":   user.ProfileId,
		})
	}
}

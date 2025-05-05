package controllers

import (
	"vnuid-identity/helpers"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type PasswordLoginRequest struct {
	DeviceId   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Email      string `json:"email" validate:"required"`
}

func LoginByPass(ctx *fiber.Ctx) error {
	var data PasswordLoginRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	// Verify google ID
	valid, user := models.VerifyPassword(data.Email, data.Password)
	if !valid {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "invalid email or password")
	}

	// Check if session already logged in on device_id
	isActive := models.CheckSession(data.DeviceId, user.ID)
	if isActive {
		// Generate token
		token, err := helpers.AddLoginSession(user, data.DeviceId, true, "pass")
		if err != nil {
			return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
		}

		return ctx.JSON(fiber.Map{"token": token})
	} else {
		var allowList []string = []string{"qr", "code", "nfc", "otp"}

		// Gen token for step 2
		token, err := utils.GenerateTemporaryToken(
			utils.TmpTokenData{
				UID:          user.ID,
				DeviceID:     data.DeviceId,
				DeviceName:   data.DeviceName,
				AllowMethods: allowList,
				Method:       "pass",
			})
		if err != nil {
			return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
		}

		return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"allow": allowList, "token": token})
	}
}

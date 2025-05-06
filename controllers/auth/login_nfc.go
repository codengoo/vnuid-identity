package controllers

import (
	"vnuid-identity/entities"
	"vnuid-identity/helpers"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type NfcLoginRequest struct {
	DeviceId   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
	NfcCode    string `json:"nfc_code" validate:"required"`
	UID        string `json:"uid" validate:"required"`
}

func LoginByNFC(ctx *fiber.Ctx) error {
	var data NfcLoginRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	// Verify google ID
	valid, user := models.VerifyNFC(data.UID, data.NfcCode)
	if !valid {
		return utils.ReturnErrorMsg(ctx, fiber.StatusUnauthorized, "invalid nfc code")
	}

	// Check if session already logged in on device_id
	isActive := models.CheckSession(data.DeviceId, user.ID)
	if isActive {
		// Generate token
		token, err := helpers.AddLoginSession(user, entities.Session{
			DeviceId:    data.DeviceId,
			DeviceName:  data.DeviceName,
			LoginMethod: "nfc",
			SavedDevice: true,
		})
		if err != nil {
			return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
		}

		return ctx.JSON(fiber.Map{"token": token, "uid": user.ProfileId})
	} else {
		var allowList []string = []string{"qr", "code", "pass", "otp"}

		// Gen token for step 2
		token, err := utils.GenerateTemporaryToken(
			utils.TmpTokenData{
				UID:          user.ID,
				DeviceID:     data.DeviceId,
				DeviceName:   data.DeviceName,
				AllowMethods: allowList,
				Method:       "nfc",
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

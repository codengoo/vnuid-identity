package controllers

import (
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type CreateNFCRequest struct {
	UserId string `json:"user_id" validate:"required"`
}

func AddNFC(ctx *fiber.Ctx) error {
	var data CreateNFCRequest

	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	if num_active := models.CountActiveNFC(data.UserId); num_active > 0 {
		return utils.ReturnErrorMsg(ctx, fiber.StatusBadRequest, "User already has 1 active NFC")
	}

	nfc_id, err := models.AddNFC(data.UserId)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"nfc_id": nfc_id})
}

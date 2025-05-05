package controllers

import (
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

func ActivateNFC(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	active := ctx.Query("active")
	err := models.SetActiveNFC(id, active == "true")
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success"})
}

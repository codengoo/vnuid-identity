package controllers

import (
	"vnuid-identity/models"

	"github.com/gofiber/fiber/v2"
)

func ActivateNFC(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	active := ctx.Query("active")
	err := models.SetActiveNFC(id, active == "true")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success"})
}

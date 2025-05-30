package controllers

import (
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type RemoveMultipleUserRequest struct {
	IDs []string `json:"ids" validate:"required"`
}

func RemoveMultipleUsers(ctx *fiber.Ctx) error {
	var data RemoveMultipleUserRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON((fiber.Map{"error": err.Error(), "msg": msg}))
	}

	if err := models.RemoveUsers(data.IDs); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "deleted records"})
}

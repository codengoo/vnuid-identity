package utils

import (
	"github.com/gofiber/fiber/v2"
)

func GetBodyData(ctx *fiber.Ctx, data interface{}) error {
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if msgs := Validate(data); msgs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid args", "msgs": msgs})
	}

	return nil
}

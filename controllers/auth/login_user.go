package controllers

import (
	"vnuid-identity/databases"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	DeviceID string `json:"device_id" validate:"required"`
}

func Login(ctx *fiber.Ctx) error {
	var data LoginRequest
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if msgs := utils.Validate(data); msgs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid args", "msgs": msgs})
	}

	var user models.User
	result := databases.DB.Where("email = ?", data.Email).First(&user) // SELECT * FROM users WHERE email = ? LIMIT 1;

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	if !utils.VerifyPassword(user.Password, data.Password) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	token, err := utils.GenerateToken(user, data.DeviceID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	return ctx.JSON(fiber.Map{"token": token})
}

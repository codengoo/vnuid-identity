package controllers

import (
	"fmt"
	"vnuid-identity/databases"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddUserRequest struct {
	Email         string `json:"email" validate:"required,email"`
	SID           string `json:"sid" validate:"required,len=8"`
	GID           string `json:"gid" validate:"required"`
	Name          string `json:"name" validate:"required"`
	OfficialClass string `json:"official_class" validate:"required"`
	Type          string `json:"type" validate:"required"`
}

func AddUser(ctx *fiber.Ctx) error {
	var data AddUserRequest

	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if msgs := utils.Validate(&data); msgs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid args", "msgs": msgs})
	}

	password, err := utils.GeneratePassword()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate password"})
	}

	user := models.User{
		Type:          data.Type,
		Email:         data.Email,
		SID:           data.SID,
		GID:           data.GID,
		Name:          data.Name,
		OfficialClass: data.OfficialClass,
		ID:            uuid.New().String(),
		Password:      password,
	}

	result := databases.DB.Create(&user)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Create user failed with message: %s", result.Error.Error()),
		})
	}

	return ctx.JSON(fiber.Map{"data": user})
}

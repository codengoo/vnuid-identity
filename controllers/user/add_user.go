package controllers

import (
	"vnuid-identity/entities"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
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
	if err := utils.GetBodyData(ctx, &data); err != nil {
		return err
	}

	user := entities.User{
		Type:          data.Type,
		Email:         data.Email,
		SID:           data.SID,
		GID:           data.GID,
		Name:          data.Name,
		OfficialClass: data.OfficialClass,
	}

	if err := models.AddUser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": user})
}

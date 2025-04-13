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
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	user, err := models.AddUser(
		entities.User{
			Type:          data.Type,
			Email:         data.Email,
			Sid:           data.SID,
			Gid:           data.GID,
			Name:          data.Name,
			OfficialClass: data.OfficialClass,
		})
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	if _, err := models.AddNFC(user.ID); err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(fiber.Map{"data": user})
}

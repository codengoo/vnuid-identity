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
	Phone         string `json:"phone" validate:"required"`
	Address       string `json:"address" validate:"required"`
	Password      string `json:"password" validate:"required"`
}

func AddUser(ctx *fiber.Ctx) error {
	var data AddUserRequest
	if err, msg := utils.GetBodyData(ctx, &data); err != nil {
		return utils.ReturnErrorDetails(ctx, fiber.StatusBadRequest, err, msg)
	}

	profile, err := models.AddUserInfo(
		entities.Profile{
			Name:          data.Name,
			Sid:           data.SID,
			Email:         data.Email,
			OfficialClass: data.OfficialClass,
			DOB:           nil,
			Phone:         data.Phone,
			Address:       data.Address,
		})
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	user, err := models.AddUser(
		entities.User{
			Type:      data.Type,
			Email:     data.Email,
			Sid:       data.SID,
			Gid:       data.GID,
			ProfileId: profile.ID,
			Password:  data.Password,
		})
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	_, err = models.AddNFC(user.ID)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.JSON(fiber.Map{"data": user})
}

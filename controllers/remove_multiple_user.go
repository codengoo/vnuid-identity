package controllers

import (
	"fmt"
	"log"
	"vnuid-identity/databases"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RemoveMultipleUserRequest struct {
	IDs []string `json:"ids" validate:"required"`
}

func isUUID(text string) bool {
	_, err := uuid.Parse(text)
	return err == nil
}

func RemoveMultipleUsers(ctx *fiber.Ctx) error {
	var data RemoveMultipleUserRequest

	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to parse request body"})
	}

	if msgs := utils.Validate(&data); msgs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid args", "msgs": msgs})
	}

	var uuids []string
	var emails []string

	for _, item := range data.IDs {
		if isUUID(item) {
			uuids = append(uuids, item)
		} else {
			emails = append(emails, item)
		}
	}

	result := databases.DB.Where("id IN ? OR email IN ?", uuids, emails).Delete(&models.User{})
	if result.Error != nil {
		log.Fatalf("could not delete records: %v", result.Error)
	}

	return ctx.JSON(fiber.Map{"message": fmt.Sprintf("Deleted %d records", result.RowsAffected)})
}

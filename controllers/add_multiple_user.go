package controllers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"vnuid-identity/databases"
	"vnuid-identity/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
)

type AddMultipleUserRequest struct {
	Email         string `xlsx:"0"`
	SID           string `xlsx:"1"`
	GID           string `xlsx:"2"`
	Name          string `xlsx:"3"`
	OfficialClass string `xlsx:"4"`
	Type          string `xlsx:"5"`
}

func AddMultipleUser(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error detach file from request"})
	}

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error opening file"})
	}
	defer fileContent.Close()

	// Read content as buffer
	fileBytes := bytes.NewBuffer(nil)
	_, err = io.Copy(fileBytes, fileContent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error reading file content"})
	}

	// create excel object to handle from this buffer
	excelFile, err := xlsx.OpenBinary(fileBytes.Bytes())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error reading file content"})
	}

	sheet := excelFile.Sheets[0]
	var users []models.User

	for i, row := range sheet.Rows {
		if i == 0 {
			continue // skip header row
		}

		var data AddMultipleUserRequest
		if err := row.ReadStruct(&data); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("error reading row %d", i)})
		}

		user := models.User{
			ID:            uuid.New().String(),
			Email:         data.Email,
			SID:           data.SID,
			GID:           data.GID,
			Name:          data.Name,
			OfficialClass: data.OfficialClass,
			Type:          data.Type,
		}
		users = append(users, user)
	}

	if result := databases.DB.Create(&users); result.Error != nil {
		log.Fatalf("could not insert batch of users: %v", result.Error)
	}

	return ctx.JSON(fiber.Map{"data": users})
}

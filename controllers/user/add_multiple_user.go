package controllers

import (
	"bytes"
	"io"
	"strconv"
	"vnuid-identity/entities"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
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

func AddMultipleUsers(ctx *fiber.Ctx) error {
	batchSize := 100

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error detach file from request"})
	}

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return utils.ReturnErrorMsg(ctx, fiber.StatusInternalServerError, "error opening file")
	}
	defer fileContent.Close()

	// Read content as buffer
	fileBytes := bytes.NewBuffer(nil)
	_, err = io.Copy(fileBytes, fileContent)
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	// create excel object to handle from this buffer
	excelFile, err := xlsx.OpenBinary(fileBytes.Bytes())
	if err != nil {
		return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
	}

	sheet := excelFile.Sheets[0]
	var users []entities.User
	for i, row := range sheet.Rows {
		if i == 0 {
			continue // skip header row
		}

		var data AddMultipleUserRequest
		if err := row.ReadStruct(&data); err != nil {
			return utils.ReturnErrorMsg(ctx, fiber.StatusInternalServerError, "error reading row: "+strconv.Itoa(i))
		}

		users = append(users,
			entities.User{
				Email:         data.Email,
				Sid:           data.SID,
				Gid:           data.GID,
				Name:          data.Name,
				OfficialClass: data.OfficialClass,
				Type:          data.Type,
			})

		if len(users) == batchSize || i == len(sheet.Rows)-1 {
			xusr, err := models.AddManyUser(users)
			if err != nil {
				return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
			}

			if _, err := models.AddManyNFC(xusr); err != nil {
				return utils.ReturnError(ctx, fiber.StatusInternalServerError, err)
			}

			users = []entities.User{}
		}
	}

	return ctx.JSON(fiber.Map{"message": "OK"})
}

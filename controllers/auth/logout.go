package controllers

import (
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
)

func Logout(ctx *fiber.Ctx) error {
	return utils.ReturnClearCookie(ctx)
}

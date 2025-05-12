package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func ReturnError(ctx *fiber.Ctx, code int, err error) error {
	return ctx.Status(code).JSON(fiber.Map{"message": "Failed", "error": err.Error()})
}

func ReturnErrorMsg(ctx *fiber.Ctx, code int, err string) error {
	return ctx.Status(code).JSON(fiber.Map{"message": "Failed", "error": err})
}

func ReturnErrorDetails(ctx *fiber.Ctx, code int, err error, msg []string) error {
	return ctx.Status(code).JSON(fiber.Map{"message": "Failed", "error": err.Error(), "msgs": msg})
}

func ReturnSuccess(ctx *fiber.Ctx, payload interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success", "data": payload})
}

func ReturnToken(ctx *fiber.Ctx, token string, uid string) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(30 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:     "uid",
		Value:    uid,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token, "uid": uid})
}

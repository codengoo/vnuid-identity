package utils

import "github.com/gofiber/fiber/v2"

func ReturnError(ctx *fiber.Ctx, code int, err error) error {
	return ctx.Status(code).JSON(fiber.Map{"error": err.Error()})
}

func ReturnErrorMsg(ctx *fiber.Ctx, code int, err string) error {
	return ctx.Status(code).JSON(fiber.Map{"error": err})
}

func ReturnErrorDetails(ctx *fiber.Ctx, code int, err error, msg []string) error {
	return ctx.Status(code).JSON(fiber.Map{"error": err.Error(), "msgs": msg})
}

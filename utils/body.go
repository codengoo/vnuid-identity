package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetBodyData(ctx *fiber.Ctx, data interface{}) (error, []string) {
	if err := ctx.BodyParser(&data); err != nil {
		return fmt.Errorf("invalid request body"), []string{}
	}

	if msgs := Validate(data); msgs != nil {
		return fmt.Errorf("invalid args"), msgs
	}

	return nil, []string{}
}

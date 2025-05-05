package middlewares

import (
	"os"
	"slices"
	"strings"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaim struct {
	utils.TokenData
	jwt.RegisteredClaims
}

func AuthCheck(roleRef []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		secretKey := os.Getenv("JWT_TOKEN")
		tokenString := ctx.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		token, err := jwt.ParseWithClaims(tokenString, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}

			role := token.Claims.(*TokenClaim).Role
			if !slices.Contains(roleRef, role) {
				return nil, fiber.ErrUnauthorized
			}

			// check more here
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		ctx.Locals("user", token.Claims)
		return ctx.Next()
	}
}

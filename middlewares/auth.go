package middlewares

import (
	"fmt"
	"os"
	"vnuid-identity/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaim struct {
	utils.TokenData
	jwt.RegisteredClaims
}

func AuthCheck(roleRef string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		secretKey := os.Getenv("JWT_TOKEN")
		tokenString := ctx.Get("Authorization")
		fmt.Println(tokenString)
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		token, err := jwt.ParseWithClaims(tokenString, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}

			role := token.Claims.(*TokenClaim).Role
			if role != roleRef {
				return nil, fiber.ErrUnauthorized
			}

			// check more here
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Invalid or expired token")
		}

		ctx.Locals("user", token.Claims)
		return ctx.Next()
	}
}

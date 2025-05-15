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
		// Lấy token từ header Authorization
		tokenString := ctx.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Nếu không có token trong header, kiểm tra cookie
		if tokenString == "" {
			tokenString = ctx.Cookies("token")
		}

		// Nếu vẫn không có token, trả về Unauthorized
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		// Phân tích JWT token
		token, err := jwt.ParseWithClaims(tokenString, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}

			// Kiểm tra role của người dùng
			role := token.Claims.(*TokenClaim).Role
			if !slices.Contains(roleRef, role) {
				return nil, fiber.ErrUnauthorized
			}

			// Kiểm tra thêm các điều kiện khác nếu cần
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		// Lưu thông tin người dùng vào context để các handler khác có thể sử dụng
		ctx.Locals("user", token.Claims)
		return ctx.Next()
	}
}

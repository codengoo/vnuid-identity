package utils

import (
	"fmt"
	"os"
	"vnuid-identity/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY string
var SECRET_KEY_2FA string

func init() {
	SECRET_KEY = os.Getenv("JWT_TOKEN")
	SECRET_KEY_2FA = os.Getenv("JWT_TOKEN_2FA")
}

type TokenData struct {
	ID       string `json:"id"`
	SID      string `json:"sid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	DeviceID string `json:"device_id"`
	Role     string `json:"role"`
}

type TmpTokenData struct {
	UID      string `json:"uid"`
	DeviceID string `json:"device_id"`
}

type TmpTokenClaim struct {
	TmpTokenData
	jwt.RegisteredClaims
}

func GenerateToken(user entities.User, deviceId string) (string, error) {
	claims := jwt.MapClaims{
		"id":        user.ID,
		"sid":       user.Sid,
		"name":      user.Name,
		"email":     user.Email,
		"device_id": deviceId,
		"role":      user.Type,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func GenerateTemporaryToken(uid string, deviceId string) (string, error) {
	claims := jwt.MapClaims{
		"uid":       uid,
		"device_id": deviceId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY_2FA))
}

func ParseTemporaryToken(tokenString string) (TmpTokenClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TmpTokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}

		// check more here
		return []byte(SECRET_KEY_2FA), nil
	})

	if err != nil || !token.Valid {
		return TmpTokenClaim{}, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(*TmpTokenClaim)
	if !ok {
		return TmpTokenClaim{}, fmt.Errorf("cannot parse token claims")
	}

	return *claims, nil
}

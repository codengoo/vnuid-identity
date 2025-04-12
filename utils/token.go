package utils

import (
	"os"
	"vnuid-identity/entities"

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
	ID       string `json:"id"`
	DeviceID string `json:"device_id"`
}

func GenerateToken(user entities.User, deviceId string) (string, error) {
	claims := jwt.MapClaims{
		"id":        user.ID,
		"sid":       user.SID,
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
		"id":        uid,
		"device_id": deviceId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY_2FA))
}

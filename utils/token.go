package utils

import (
	"os"
	"vnuid-identity/models"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY string

func init() {
	SECRET_KEY = os.Getenv("JWT_TOKEN")
}

type TokenData struct {
	SID      string `json:"sid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	DeviceID string `json:"device_id"`
	Role     string `json:"role"`
}

func GenerateToken(user models.User, deviceId string) (string, error) {
	claims := jwt.MapClaims{
		"sid":       user.SID,
		"name":      user.Name,
		"email":     user.Email,
		"device_id": deviceId,
		"role":      user.Type,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

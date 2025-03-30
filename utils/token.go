package utils

import (
	"vnuid-identity/models"

	"github.com/golang-jwt/jwt/v5"
)

const SECRET_KEY = "your_secret_key"

func GenerateToken(user models.User, deviceId string) (string, error) {
	claims := jwt.MapClaims{
		"sid":       user.SID,
		"name":      user.Name,
		"email":     user.Email,
		"device_id": deviceId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

package utils

import (
	"fmt"
	"os"
	"vnuid-identity/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenData struct {
	UID      string `json:"uid"`
	SID      string `json:"sid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	DeviceID string `json:"device_id"`
	Role     string `json:"role"`
}

type TmpTokenData struct {
	UID          string   `json:"uid"`
	DeviceID     string   `json:"device_id"`
	DeviceName   string   `json:"device_name"`
	AllowMethods []string `json:"allow_methods"`
	Method       string   `json:"method"`
}

type QRTokenData struct {
	UID      string `json:"uid"`
	DeviceID string `json:"device_id"`
}

type TmpTokenClaim struct {
	TmpTokenData
	jwt.RegisteredClaims
}

type QRTokenClaim struct {
	QRTokenData
	jwt.RegisteredClaims
}

func GenerateToken(user entities.User, deviceId string) (string, error) {
	var SECRET_KEY = os.Getenv("JWT_TOKEN")
	claims := jwt.MapClaims{
		"uid": user.ProfileId,
		"sid": user.Sid,
		// "name":      user.Name,
		"email":     user.Email,
		"device_id": deviceId,
		"role":      user.Type,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func GenerateTemporaryToken(data TmpTokenData) (string, error) {
	var SECRET_KEY_2FA = os.Getenv("JWT_TOKEN_2FA")
	claims := jwt.MapClaims{
		"uid":           data.UID,
		"device_id":     data.DeviceID,
		"device_name":   data.DeviceName,
		"allow_methods": data.AllowMethods,
		"method":        data.Method,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY_2FA))
}

func ParseTemporaryToken(tokenString string) (TmpTokenClaim, error) {
	var SECRET_KEY_2FA = os.Getenv("JWT_TOKEN_2FA")
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

func PrintAllClaims(claims jwt.Claims) {
	for k, v := range claims.(jwt.MapClaims) {
		fmt.Println(k, v)
	}
}

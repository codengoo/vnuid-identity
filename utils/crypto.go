package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func EncryptPayload(key string, uid string, device_id string) (string, error) {
	aesKey := []byte(key)
	timestamp := time.Now().UnixMilli()
	payload := fmt.Sprintf("%s:%s:%d", uid, device_id, timestamp)
	plaintext := []byte(payload)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptPayload(key string, encrypted string) (uid, deviceID string, err error) {
	aesKey := []byte(key)
	cipherText, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", "", errors.New("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	parts := strings.Split(string(cipherText), ":")
	if len(parts) != 3 {
		return "", "", errors.New("invalid payload format")
	}

	uid = parts[0]
	deviceID = parts[1]
	tsMillis, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", "", errors.New("invalid timestamp")
	}

	now := time.Now().UnixMilli()
	delta := now - tsMillis
	if delta < 0 {
		delta = -delta
	}

	if delta > 2000 {
		return "", "", errors.New("token expired or time drift too large")
	}

	return uid, deviceID, nil
}

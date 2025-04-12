package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"vnuid-identity/databases"
	"vnuid-identity/models"
	"vnuid-identity/utils"

	"github.com/gofiber/websocket/v2"
)

func genQRAccept(uid string, deviceId string, save bool) (string, error) {
	user, err := models.GetUser(uid)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	token, err := utils.GenerateToken(user, deviceId)
	if err != nil {
		return "", fmt.Errorf("error generating token")
	}

	// Create login session
	if _, err := models.CreateSession(deviceId, user.ID, save); err != nil {
		return "", fmt.Errorf("create session: %s", err.Error())
	}

	return token, nil
}

func sendMessage(ctx *websocket.Conn, text string) {
	err := ctx.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func ListenLogin2FA(ctx *websocket.Conn) {
	defer func() {
		fmt.Println("Client disconnected")
		ctx.Close()
	}()

	session := ctx.Params("session")

	bgctx := context.Background()
	_, err := databases.RD.Get(bgctx, fmt.Sprintf("%s%s", LOGIN_KEY, session)).Result()

	if err != nil {
		sendMessage(ctx, err.Error())
		ctx.Close()
		return
	}

	for {
		val, err := databases.RD.Get(bgctx, fmt.Sprintf("%s%s", LOGIN_ACCEPT_KEY, session)).Result()

		if err == nil {
			var content LoginByQr2FaAcceptConfig
			err = json.Unmarshal([]byte(val), &content)
			if err != nil {
				sendMessage(ctx, err.Error())
			} else {
				token, err := genQRAccept(content.UID, content.DeviceID, content.SaveDevice)
				if err != nil {
					sendMessage(ctx, err.Error())
				}

				sendMessage(ctx, fmt.Sprintf("token::%s", token))
				return
			}
		}
	}
}

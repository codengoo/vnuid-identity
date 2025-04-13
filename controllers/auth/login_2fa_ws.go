package controllers

import (
	"fmt"
	"vnuid-identity/helpers"
	"vnuid-identity/models"

	"github.com/gofiber/websocket/v2"
)

func genToken(uid string, deviceId string, save bool) (string, error) {
	user, err := models.GetUser(uid)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	return helpers.AddLoginSession(user, deviceId, save)
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
	_, err := models.GetLoginSession(session)
	if err != nil {
		sendMessage(ctx, err.Error())
		ctx.Close()
		return
	}

	for {
		content, err := models.GetLoginAcceptSession(session)
		if err != nil {
			sendMessage(ctx, err.Error())
		} else {
			token, err := genToken(content.UID, content.DeviceID, content.SaveDevice)
			if err != nil {
				sendMessage(ctx, err.Error())
			}

			sendMessage(ctx, fmt.Sprintf("token::%s", token))
			return
		}
	}
}

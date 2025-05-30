package controllers

import (
	"fmt"
	"vnuid-identity/entities"
	"vnuid-identity/helpers"
	"vnuid-identity/models"

	"github.com/gofiber/websocket/v2"
)

func genToken(uid string, deviceId string, save bool, method string) (string, error) {
	user, err := models.GetUser(uid)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	return helpers.AddLoginSession(user, entities.Session{
		DeviceId:    deviceId,
		DeviceName:  "",
		LoginMethod: method,
		SavedDevice: save,
	})
}

func sendMessage(ctx *websocket.Conn, text string) {
	err := ctx.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func ListenLogin(ctx *websocket.Conn) {
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
		if err == nil {
			token, err := genToken(content.UID, content.DeviceID, content.SaveDevice, content.Method)
			if err != nil {
				sendMessage(ctx, err.Error())
			}

			sendMessage(ctx, fmt.Sprintf("token::%s", token))
			return
		}
	}
}

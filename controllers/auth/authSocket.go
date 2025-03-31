package controllers

import (
	"fmt"

	"github.com/gofiber/websocket/v2"
)

func AuthSocket(ctx *websocket.Conn) {
	defer ctx.Close()

	fmt.Println("Client đã kết nối!")
}

package server

import (
	"io"

	"golang.org/x/net/websocket"
)

func ListenHandler(ws *websocket.Conn) {
	logger.Printf("web socket connection!!")
	io.Copy(ws, ws)
}

package ws

import "github.com/gorilla/websocket"

const (
	// SocketBufferSize is the buffer
	SocketBufferSize = 1024

	// MessageBufferSize is the buffer
	MessageBufferSize = 256
)

// Upgrader is the websocket upgrader
var Upgrader = &websocket.Upgrader{
	ReadBufferSize:  SocketBufferSize,
	WriteBufferSize: SocketBufferSize,
}

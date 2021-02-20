package ws

import (
	"github.com/gorilla/websocket"
)

const (
	// SocketBufferSize is the buffer
	SocketBufferSize = 1024
	// MessageBufferSize is the buffer
	MessageBufferSize = 256
)

// Client represents a single websocket connection.
type Client struct {
	// represents the socket connection for a single chat client.
	socket *websocket.Conn
	// send is a buffered channel through which messages are sent
	// to the socket connection.
	send chan *Message
	// room represents the room this client is in.
	room *Room
}

// NewClient does what you'd expect right?
func NewClient() *Client {

}

// Read an incoming socket message
func (c *Client) Read() {

}

// Write a message to the connection
func (c *Client) Write() {

}

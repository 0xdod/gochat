package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

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

// Client represents a single websocket connection.
type Client struct {
	id string
	// represents the socket connection for a single chat client.
	socket *websocket.Conn
	// send is a buffered channel through which messages are sent
	// to the socket connection.
	send chan *Message
	// room represents the room this client is in.
	room *Room
}

// NewClient does what you'd expect right?
func NewClient(s *websocket.Conn) *Client {
	return &Client{
		id:     "client",
		socket: s,
		send:   make(chan *Message, MessageBufferSize),
	}
}

// Read an incoming message from the socket connection.
func (c *Client) Read() {
	for {
		var msg = &Message{
			From: c.id,
			When: time.Now(),
		}
		if err := c.socket.ReadJSON(msg); err != nil {
			break
		}
		c.room.forward <- msg
	}
	c.socket.Close()
}

// Write a message to the connection.
func (c *Client) Write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

func (c *Client) String() string {
	return c.id
}

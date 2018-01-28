package chat

//Defines a client using the chat app
import (
	"github.com/gorilla/websocket"
)

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
)

//client represents a single chatting user
type Client struct {
	//represents a socket connection for a single user
	socket *websocket.Conn
	//send is a buffered channel through which messages are sent
	send chan *Message
	//room is the room this client is chatting in
	room *Room
	// data finder knows how to find client information
	DataFinder
	IsPresent bool
}

func NewClient(s *websocket.Conn, df DataFinder, ip bool) *Client {
	return &Client{
		socket:     s,
		send:       make(chan *Message, MessageBufferSize),
		DataFinder: df,
		IsPresent:  ip,
	}
}

func (c *Client) Read() {
	defer c.socket.Close()
	for {
		msg := NewMessage(nil, c.DataFinder, "")
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *Client) Write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}

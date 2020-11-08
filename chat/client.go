package chat

//Defines a client using the chat app
import (
	"time"

	"github.com/fibreactive/chat/models"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	SocketBufferSize  = 1024
	MessageBufferSize = 256
)

//client represents a single chatting user
type Client struct {
	id string
	//represents a socket connection for a single user
	socket *websocket.Conn
	//send is a buffered channel through which messages are sent
	send chan *message
	//room is the room this client is chatting in
	room *Room
	// userData holds information about the user
	user *models.UserModel
}

func NewClient(s *websocket.Conn, user *models.UserModel) *Client {
	return &Client{
		id:     uuid.Must(uuid.NewV4()).String(),
		socket: s,
		send:   make(chan *message, MessageBufferSize),
		user:   user,
	}
}

// export this one
func (c *Client) Read() {
	defer c.socket.Close()
	for {
		msg := &message{}
		err := c.socket.ReadJSON(msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.From = c.user.Nickname
		if c.user.AvatarURL != "" {
			// check might not be needed
			msg.AvatarURL = c.user.AvatarURL
		}
		c.room.forward <- msg
	}
}

// export this one too
func (c *Client) Write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}

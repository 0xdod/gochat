package chat

//Defines a client using the chat app
import (
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/errors/fmt"
)

//type message is a representation of messages sent between users
type message struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	When      time.Time `json:"when"`
	AvatarURL string    `json:"avatar_url"`
}

//client represents a single chatting user
type client struct {
	//represents a socket connection for a single user
	socket *websocket.Conn
	//send is a buffered channel through which messages are sent
	send chan *message
	//room is the room this client is chatting in
	room *room
	// userData holds information about the user
	userData map[string]interface{}
}

// export this one
func (c *client) read() {
	defer c.socket.Close()
	for {
		msg := &message{}
		err := c.socket.ReadJSON(msg)
		if err != nil {
			fmt.Println("")
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.UserID = c.userData["userid"].(string)
		if avatarUrl, ok := c.userData["avatar_url"]; ok {
			// check might not be needed
			msg.AvatarURL = avatarUrl.(string)
		}
		c.room.forward <- msg
	}
}

// export this one too
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}

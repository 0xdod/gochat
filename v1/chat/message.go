package chat

import (
	"fmt"
	"time"
)

const (
	NEW_USER    = "welcome"
	USER_LEFT   = "leave"
	USER_JOINED = "join"
)

//type message is a representation of messages sent between users
type Message struct {
	From      string    `json:"from,omitempty"`
	To        string    `json:"to,omitempty"`
	Message   string    `json:"message,omitempty"`
	When      time.Time `json:"when,omitempty"`
	AvatarURL string    `json:"avatarURL,omitempty"`
	UserID    int       `json:"userID,omitempty"`
}

func (m *Message) String() string {
	return fmt.Sprintf("From: %s, To: %s, Message: %s", m.From, m.To, m.Message)
}

func NewMessage(to, from DataFinder, content string) *Message {
	msg := &Message{}
	if from == nil {
		msg.From = "Admin"
	} else {
		msg.From = from.GetName()
		msg.UserID = from.GetIntID()
		msg.AvatarURL = from.GetAvatarURL()
	}
	if to != nil {
		msg.To = to.GetName()
	}
	msg.Message = content
	msg.When = time.Now()
	return msg
}

func generateAdminMessage(c *Client, info string) *Message {
	var msg string
	roomName := c.room.GetName()
	username := c.GetName()
	switch info {
	case NEW_USER:
		msg = "/Hello " + username + ", Welcome to the " + roomName + " chat room."
	case USER_JOINED:
		msg = "/" + username + " has joined!."
	case USER_LEFT:
		msg = "/" + username + " has left!."
	}
	return NewMessage(nil, nil, msg)
}

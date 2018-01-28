package chat

import (
	"time"
)

const (
	NEW_USER    = "welcome"
	USER_LEFT   = "leave"
	USER_JOINED = "join"
)

//type message is a representation of messages sent between users
type message struct {
	From      string    `json:"from,omitempty"`
	To        string    `json:"to,omitempty"`
	Message   string    `json:"message,omitempty"`
	When      time.Time `json:"when,omitempty"`
	AvatarURL string    `json:"avatarURL,omitempty"`
	UserID    uint      `json:"userID,omitempty"`
}

func NewMessage(to, from *Client, content string) *message {
	msg := &message{}
	if from == nil {
		msg.From = "Admin"
	} else {
		msg.From = from.user.Nickname
		msg.UserID = from.user.ID
		msg.AvatarURL = from.user.AvatarURL
	}
	if to != nil {
		msg.To = to.user.Nickname
	}
	msg.Message = content
	msg.When = time.Now()
	return msg
}

func generateAdminMessage(c *Client, info string) *message {
	var msg string
	var roomName string
	username := c.user.Nickname
	if c.room.room == nil {
		roomName = "default"
	} else {
		roomName = c.room.room.Name
	}
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

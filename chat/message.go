package chat

import "time"

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
	ClientID  string    `json:"clientID,omitempty"`
}

func NewMessage(to *Client, from, body string) *message {
	msg := &message{
		From:    from,
		Message: body,
		When:    time.Now(),
	}
	if to != nil {
		msg.To = to.user.Nickname
		msg.ClientID = to.id
	}
	return msg
}

func generateAdminMessage(c *Client, info string) *message {
	var msg string
	username := c.user.Nickname
	roomName := c.room.name
	switch info {
	case NEW_USER:
		msg = "Hello " + username + ", Welcome to the " + roomName + " chat room."
	case USER_JOINED:
		msg = username + " has joined!."
	case USER_LEFT:
		msg = username + " has left!."
	}
	return NewMessage(nil, "Admin", msg)
}

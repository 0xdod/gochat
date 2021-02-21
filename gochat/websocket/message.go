package ws

import (
	"fmt"
	"time"
)

const (
	NewUser    = "welcome"
	UserLeft   = "leave"
	UserJoined = "join"
)

//type message is a representation of messages sent between users
type Message struct {
	From    string    `json:"from,omitempty"`
	Message string    `json:"message,omitempty"`
	When    time.Time `json:"when,omitempty"`
	UserID  int       `json:"userID,omitempty"`
}

func (m *Message) String() string {
	return fmt.Sprintf("From: %s\n, Message: %s\n At: %v\n", m.From, m.Message, m.When)
}

package gochat

import "context"

// Message represents chat messages exchanged in a room.
type Message struct {
	Model
	RoomID int    `json:"room_id" gorm:"not null"`
	UserID int    `json:"user_id" gorm:"not null"`
	Text   string `json:"text,omitempty" gorm:"not null"`
	Edited bool   `json:"edited,omitempty"`
}

// MessageService represents a service for managing messages.
type MessageService interface {
	// Retrieves a Message by ID.
	// Returns ENOTFOUND if Message does not exist.
	FindMessageByID(ctx context.Context, id int) (*Message, error)

	// Retrieves a list of Messages by filter. Also returns total count of matching
	// Messages which may differ from returned results if filter.Limit is specified.
	FindMessages(ctx context.Context, filter MessageFilter) ([]*Message, int, error)

	// Creates a new Message.
	CreateMessage(ctx context.Context, Message *Message) error

	// Updates a Message object.
	UpdateMessage(ctx context.Context, id int, upd MessageUpdate) (*Message, error)

	// Permanently deletes a Message.
	DeleteMessage(ctx context.Context, id int) error
}

// MessageFilter represents a filter passed to FindMessages().
type MessageFilter struct {
	// Filtering fields.
	ID *int `json:"id"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// MessageUpdate represents a set of fields to be updated via UpdateMessage().
type MessageUpdate struct {
	Text *string `json:"text"`
}

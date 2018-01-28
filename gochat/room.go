package gochat

import (
	"context"
)

// Room represents a single chat room.
type Room struct {
	Model
	Name         string     `json:"name" gorm:"not null,size:100"`
	Description  *string    `json:"description,omitempty"`
	Link         string     `json:"link" gorm:"not null"`
	Invite       string     `json:"invite" gorm:"not null,uniqueIndex"`
	Icon         *string    `json:"icon"`
	Participants []*User    `json:"participants" gorm:"many2many:room_participants;"`
	Admins       []*User    `json:"admins" gorm:"many2many:room_admins;"`
	Messages     []*Message `gorm:"constraint:OnDelete:CASCADE;"`
}

func (r *Room) InviteLink() string {
	return ""
}

// RoomService represents a service for managing messages.
type RoomService interface {
	// Retrieves a Room by ID.
	// Returns ENOTFOUND if Room does not exist.
	FindRoomByID(ctx context.Context, id int) (*Room, error)

	// Retrieves a list of Rooms by filter. Also returns total count of matching
	// Rooms which may differ from returned results if filter.Limit is specified.
	FindRooms(ctx context.Context, filter RoomFilter) ([]*Room, int, error)

	// Creates a new Room.
	CreateRoom(ctx context.Context, Room *Room) error

	// Updates a Room object.
	UpdateRoom(ctx context.Context, id int, upd RoomUpdate) (*Room, error)

	// Permanently deletes a Room.
	DeleteRoom(ctx context.Context, id int) error
}

// RoomFilter represents a filter passed to FindRooms().
type RoomFilter struct {
	// Filtering fields.
	ID   *int   `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// RoomUpdate represents a set of fields to be updated via UpdateRoom().
type RoomUpdate struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Link        string  `json:"link"`
	Icon        string  `json:"icon"`
}

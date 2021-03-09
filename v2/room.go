package gochat

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Room represents a single chat room.
type Room struct {
	Model
	Name         string     `json:"name" gorm:"not null,size:100"`
	Description  *string    `json:"description,omitempty"`
	Link         string     `json:"link" gorm:"not null"`
	Invite       string     `json:"invite" gorm:"not null,uniqueIndex"`
	Icon         *string    `json:"icon"`
	CreatorID    int        `json:"creator_id" gorm:"not null"`
	Participants []*User    `json:"participants" gorm:"many2many:room_participants;"`
	Admins       []*User    `json:"admins" gorm:"many2many:room_admins;"`
	Messages     []*Message `gorm:"constraint:OnDelete:CASCADE;"`
}

var alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomSequence() string {
	var seq string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		randIndex := rand.Intn(len(alphabets))
		seq += string(alphabets[randIndex])
	}
	return seq
}

func (r *Room) InviteLink(domain string) string {
	r.Invite = randomSequence()
	r.Link = fmt.Sprintf("http://%s/chat/%s", domain, r.Invite)
	return r.Link
}

func (r *Room) AddParticipant(user *User) *Room {
	r.Participants = append(r.Participants, user)
	return r
}

func (r *Room) AddAdmin(user *User) *Room {
	r.Admins = append(r.Admins, user)
	return r
}

func (r *Room) IsParticipant(u *User) bool {
	for _, v := range r.Participants {
		if v.ID == u.ID {
			return true
		}
	}
	return false
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
	ID     *int    `json:"id"`
	Name   *string `json:"name"`
	Link   *string `json:"link"`
	Invite *string `json:"invite"`
	UserID *int    `json:"user_id"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type Filter map[string]interface{}

// RoomUpdate represents a set of fields to be updated via UpdateRoom().
type RoomUpdate struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	Link         *string `json:"link"`
	Icon         *string `json:"icon"`
	Participants []*User `json:"participants"`
	Admins       []*User
}

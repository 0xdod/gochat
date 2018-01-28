package gochat

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Model adds basic database fields.
type Model struct {
	ID        int        `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

// User represents a user in the system. Users are typically created after
// signing up via username and password.
type User struct {
	Model
	Name      string     `json:"name" gorm:"size:255,not null"`
	Username  string     `json:"username" gorm:"size:50,uniqueIndex,not null"`
	Email     string     `json:"email" gorm:"size:255,uniqueIndex,not null"`
	Password  string     `json:"password" gorm:"not null,size:255"`
	AvatarURL *string    `json:"avatar_url"`
	Messages  []*Message `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Rooms     []*Room    `gorm:"many2many:room_participants;"`
}

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil

}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UserService represents a service for managing users.
type UserService interface {
	// Retrieves a user by ID.
	// Returns ENOTFOUND if user does not exist.
	FindUserByID(ctx context.Context, id int) (*User, error)

	// Retrieves a list of users by filter. Also returns total count of matching
	// users which may differ from returned results if filter.Limit is specified.
	FindUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)

	// Creates a new user.
	CreateUser(ctx context.Context, user *User) error

	// Updates a user object. Returns EUNAUTHORIZED if current user is not
	// the user that is being updated. Returns ENOTFOUND if user does not exist.
	UpdateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)

	// Permanently deletes a user and all owned dials. Returns EUNAUTHORIZED
	// if current user is not the user being deleted. Returns ENOTFOUND if
	// user does not exist.
	DeleteUser(ctx context.Context, id int) error

	// Authenticate finds a user with details provided.
	Authenticate(ctx context.Context, email, password string) *User
}

// UserFilter represents a filter passed to FindUsers().
type UserFilter struct {
	// Filtering fields.
	ID       *int    `json:"id"`
	Email    *string `json:"email"`
	Username *string `json:"username"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

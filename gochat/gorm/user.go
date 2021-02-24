package gorm

import (
	"context"

	"github.com/0xdod/gochat/gochat"
	"gorm.io/gorm"
)

// UserGorm implements the gochat.UserService interface
type UserGorm struct {
	*gorm.DB
}

func NewUserService(db *DB) *UserGorm {
	return &UserGorm{db.DB}
}

func (us *UserGorm) FindUserByID(ctx context.Context, id int) (*gochat.User, error) {
	return nil, nil
}

func (us *UserGorm) FindUsers(ctx context.Context, filter gochat.UserFilter) ([]*gochat.User, int, error) {
	return nil, 0, nil
}

func (us *UserGorm) UpdateUser(ctx context.Context, id int, upd gochat.UserUpdate) (*gochat.User, error) {
	return nil, nil
}

func (us *UserGorm) DeleteUser(ctx context.Context, id int) error {
	return nil
}

func (us *UserGorm) CreateUser(ctx context.Context, user *gochat.User) error {
	return nil
}

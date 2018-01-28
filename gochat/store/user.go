package store

import (
	"context"

	"github.com/0xdod/gochat/gochat"
	"gorm.io/gorm"
)

type M map[string]interface{}

// UserGorm implements the gochat.UserService interface
type UserGorm struct {
	*gorm.DB
}

func NewUserStore(db *DB) *UserGorm {
	return &UserGorm{db.DB}
}

func (us *UserGorm) FindUserByID(ctx context.Context, id int) (*gochat.User, error) {
	user := &gochat.User{}
	err := us.DB.WithContext(ctx).First(user, id).Error

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserGorm) Authenticate(ctx context.Context, email, password string) *gochat.User {

	user := &gochat.User{}
	err := us.DB.WithContext(ctx).Where("email = ?", email).First(user).Error

	if err != nil {
		return nil
	}
	if user.VerifyPassword(password) {
		return user
	}
	return nil

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
	err := us.DB.WithContext(ctx).Create(user).Error
	if isDuplicateKeyError(err) {
		return gochat.ECONFLICT
	}
	return err
}

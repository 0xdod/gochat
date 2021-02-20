package postgres

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	User *gochat.User
}

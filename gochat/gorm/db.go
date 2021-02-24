package gorm

import (
	"github.com/0xdod/gochat/gochat"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func ConnectToDB(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&gochat.User{}, &gochat.Room{})
	return &DB{db}, err
}

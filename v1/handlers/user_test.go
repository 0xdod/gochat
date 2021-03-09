package handlers

import (
	"github.com/0xdod/gochat/models"
	"github.com/jinzhu/gorm"
)

var testUsers = []*models.User{
	&models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "Damilola",
	},
	&models.User{
		Model: gorm.Model{
			ID: 2,
		},
		Username: "Timileyin",
	},
	&models.User{
		Model: gorm.Model{
			ID: 3,
		},
		Username: "Grace",
	},
	&models.User{
		Model: gorm.Model{
			ID: 4,
		},
		Username: "Pelumi",
	},
}

type TestUserService struct{}

func (tus *TestUserService) Authenticate(email, password string) *models.User {
	return nil
}

func (tus *TestUserService) FindByID(id uint) *models.User {
	for _, u := range testUsers {
		if id == u.Model.ID {
			return u
		}
	}
	return nil
}

func (tus *TestUserService) FindByEmail(email string) *models.User {
	return nil
}

func (tus *TestUserService) Create(user *models.User) error {
	return nil
}

func (tus *TestUserService) Update(user *models.User) error {
	return nil
}

func (tus *TestUserService) GetRooms(user *models.User) []*models.Room {
	return nil
}

func (tus *TestUserService) Delete(id uint) error {
	return nil
}

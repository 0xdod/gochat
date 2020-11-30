package handlers

import (
	"github.com/fibreactive/chat/models"
	"github.com/jinzhu/gorm"
)

var testRooms = []*models.Room{
	&models.Room{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "beasts",
	},
	&models.Room{
		Model: gorm.Model{
			ID: 2,
		},
	},
	&models.Room{
		Model: gorm.Model{
			ID: 3,
		},
	},
	&models.Room{
		Model: gorm.Model{
			ID: 4,
		},
	},
}

type TestRoomService struct{}

func (trs *TestRoomService) GetAll() []*models.Room {
	return nil
}

func (trs *TestRoomService) FindByID(id uint) *models.Room {
	return nil
}

func (trs *TestRoomService) FindByLink(l string) *models.Room {
	return nil
}

func (trs *TestRoomService) FindMany(m models.Map) []*models.Room {
	return nil
}

func (trs *TestRoomService) GetAdmins(r *models.Room) []*models.User {
	return nil
}

func (trs *TestRoomService) GetParticipants(r *models.Room) []*models.User {
	return nil
}

func (trs *TestRoomService) GetMessages(r *models.Room) []*models.Message {
	return nil
}

func (trs *TestRoomService) AddParticipant(r *models.Room, u *models.User) error {
	return nil
}

func (trs *TestRoomService) RemoveParticipant(r *models.Room, u *models.User) error {
	return nil
}

func (trs *TestRoomService) Create(m *models.Room) error {
	return nil
}

func (trs *TestRoomService) Update(m *models.Room) error {
	return nil
}

func (trs *TestRoomService) Delete(id uint) error {
	return nil
}

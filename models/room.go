package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type RoomModel struct {
	gorm.Model
	Name        string `gorm:"not null;unique_index"`
	Description string
	Users       []UserModel `gorm:"many2many:room_participants"`
}

type RoomService interface {
	All() []*RoomModel
	FindByID(id uint) *RoomModel
	FindByName(name string) *RoomModel
	Create(room *RoomModel) error
	Update(room *RoomModel) error
	Delete(id uint) error
}

type RoomGorm struct {
	*gorm.DB
}

func NewRoomGorm(connInfo string) (*RoomGorm, error) {
	db, err := gorm.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}
	return &RoomGorm{db}, nil
}

func (rg *RoomGorm) All() []RoomModel {
	var rm []RoomModel
	rg.DB.Find(&rm)
	return rm
}

func (rg *RoomGorm) FindByID(id uint) *RoomModel {
	return rg.byQuery(rg.DB.Where("id = ?", id))
}

func (rg *RoomGorm) FindByName(name string) *RoomModel {
	return rg.byQuery(rg.DB.Where("name = ?", name))
}

func (rg *RoomGorm) Create(room *RoomModel) error {
	return rg.DB.Create(room).Error
}

func (rg *RoomGorm) Update(room *RoomModel) error {
	return rg.DB.Save(room).Error
}

func (rg *RoomGorm) Delete(id uint) error {
	room := rg.FindByID(id)
	return rg.DB.Delete(room).Error
}

func (rg *RoomGorm) byQuery(query *gorm.DB) *RoomModel {
	room := &RoomModel{}
	err := query.First(room).Error
	switch err {
	case nil:
		return room
	case gorm.ErrRecordNotFound:
		return nil
	default:
		panic(err)
	}
}

func (rg *RoomGorm) DestructiveReset() {
	rg.DropTableIfExists(&RoomModel{})
	rg.AutoMigrate()
}

func (rg *RoomGorm) AutoMigrate() {
	rg.DB.AutoMigrate(&RoomModel{})
}

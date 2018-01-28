package models

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Room struct {
	gorm.Model
	Name        string `gorm:"not null;"`
	Description string
	Link        string     `gorm:"not null;unique_index"`
	Users       []*User    `gorm:"many2many:room_participants"`
	Admins      []*User    `gorm:"many2many:room_admins"`
	Messages    []*Message `gorm:"constraint:OnDelete:CASCADE;"`
	AvatarURL   string
}

func (m Room) String() string {
	return fmt.Sprintf("%s", m.Name)
}

func (rm *Room) GetIntID() int {
	return int(rm.ID)
}

func (rm *Room) GetStringID() string {
	return ""
}

func (rm *Room) GetName() string {
	return rm.Name
}

func (rm *Room) GetAvatarURL() string {
	return rm.AvatarURL
}

func (rm *Room) IsNewClient() bool {
	return true
}

func (rm *Room) CreateLink() {
	id := strconv.Itoa(int(rm.ID))
	buf := &bytes.Buffer{}
	wc := base64.NewEncoder(base64.URLEncoding, buf)
	wc.Write([]byte(id + ":" + rm.Name))
	defer wc.Close()
	rm.Link = buf.String()
}

func (rm *Room) AfterCreate(tx *gorm.DB) (err error) {
	rm.CreateLink()
	return tx.Model(rm).Update("link", rm.Link).Error
}

func RoomIDFromLink(s string) int {
	reader := strings.NewReader(s)
	decoder := base64.NewDecoder(base64.URLEncoding, reader)
	p := make([]byte, reader.Len())
	decoder.Read(p)
	s = string(p)
	ss := strings.Split(s, ":")
	id, _ := strconv.Atoi(ss[0])
	return id
}

type Map map[string]interface{}

type RoomService interface {
	GetAll() []*Room
	FindMany(Map) []*Room
	FindByID(id uint) *Room
	FindByLink(string) *Room
	Create(room *Room) error
	Update(room *Room) error
	Delete(id uint) error
	GetAdmins(*Room) []*User
	GetParticipants(*Room) []*User
	AddParticipant(*Room, *User) error
	RemoveParticipant(*Room, *User) error
	GetMessages(r *Room) []*Message
}

type RoomGorm struct {
	*gorm.DB
}

func NewRoomGorm(db *gorm.DB) *RoomGorm {
	return &RoomGorm{db}
}

func (rg *RoomGorm) GetAll() []*Room {
	var rm []*Room
	rg.DB.Find(&rm).Order("created_at desc")
	return rm
}

func (rg *RoomGorm) FindByLink(s string) *Room {
	id := RoomIDFromLink(s)
	return rg.FindByID(uint(id))
}

func (rg *RoomGorm) FindByID(id uint) *Room {
	return rg.byQuery(rg.DB.Where("id = ?", id))
}

func (rg *RoomGorm) FindMany(m Map) []*Room {
	var rm []*Room
	rg.DB.Where(m).Find(&rm)
	return rm
}

func (rg *RoomGorm) Create(room *Room) error {
	return rg.DB.Create(room).Error
}

func (rg *RoomGorm) Update(room *Room) error {
	return rg.DB.Save(room).Error
}

func (rg *RoomGorm) Delete(id uint) error {
	room := Room{}
	return rg.DB.Delete(&room, id).Error
}

func (rg *RoomGorm) GetParticipants(r *Room) []*User {
	var users []*User
	rg.DB.Model(r).Association("Users").Find(&users)
	return users
}

func (rg *RoomGorm) IsUserPresent(r *Room, u *User) bool {
	users := rg.GetParticipants(r)
	for _, user := range users {
		if user.ID == u.ID {
			return true
		}
	}
	return false
}

func (rg *RoomGorm) AddParticipant(room *Room, u *User) error {
	return rg.DB.Model(room).Association("Users").Append(u).Error
}

func (rg *RoomGorm) RemoveParticipant(room *Room, u *User) error {
	return rg.DB.Model(room).Association("Users").Delete(u).Error
}

func (rg *RoomGorm) GetAdmins(r *Room) []*User {
	var users []*User
	rg.DB.Model(r).Association("Admins").Find(&users)
	return users
}

func (rg *RoomGorm) GetMessages(r *Room) []*Message {
	var messages []*Message
	rg.DB.Model(r).Order("created_at").Association("Messages").Find(&messages)
	return messages
}

func (rg *RoomGorm) DestructiveReset() {
	rg.DropTableIfExists(&Room{})
	rg.AutoMigrate()
}

func (rg *RoomGorm) AutoMigrate() {
	rg.DB.AutoMigrate(&Room{})
}

func (rg *RoomGorm) byQuery(query *gorm.DB) *Room {
	room := &Room{}
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

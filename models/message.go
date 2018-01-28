package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	Message string
	UserID  uint
	User    *User
	RoomID  uint
	Room    *Room `gorm:"constraint:OnDelete:CASCADE;"`
}

func (m Message) String() string {
	return fmt.Sprintf("%s", m.Message)
}

type MessageService interface {
	GetAll() []*Message
	GetByRoom(*Room) []*Message
	GetByUser(*User) []*Message
	FindOne(Map) *Message
	FindMany(Map) []*Message
	Create(*Message) error
	Update(*Message) error
	Delete(uint) error
}

type MessageGorm struct {
	*gorm.DB
}

func NewMessageGorm(db *gorm.DB) *MessageGorm {
	return &MessageGorm{db}
}

// will create new message based on message
// get user, get room, get message

func (mg *MessageGorm) Create(m *Message) error {
	return mg.DB.Create(m).Error
}

func (mg *MessageGorm) Update(m *Message) error {
	return mg.DB.Save(m).Error
}

func (mg *MessageGorm) Delete(id uint) error {
	var m Message
	return mg.DB.Delete(&m, id).Error
}

func (mg *MessageGorm) GetAll() []*Message {
	var msgs []*Message
	mg.DB.Find(&msgs)
	return msgs
}

func (mg *MessageGorm) GetByUser(u *User) []*Message {
	var msgs []*Message
	mg.Model(u).Order("created_at desc").Association("Messages").Find(msgs)
	return msgs
}

func (mg *MessageGorm) GetByRoom(r *Room) []*Message {
	var msgs []*Message
	mg.Model(r).Order("created_at desc").Association("Messages").Find(msgs)
	return msgs
}

func (mg *MessageGorm) FindOne(m Map) *Message {
	var msg *Message
	mg.DB.First(&msg, m)
	return msg
}

func (mg *MessageGorm) FindMany(m Map) []*Message {
	var msgs []*Message
	mg.DB.Find(msgs, m).Order("created_at desc")
	return msgs
}

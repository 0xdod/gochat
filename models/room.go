package models

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type RoomModel struct {
	gorm.Model
	Name        string `gorm:"not null;"`
	Description string
	Link        string       `gorm:"not null;unique_index"`
	Users       []*UserModel `gorm:"many2many:room_participants"`
	Admins      []*UserModel `gorm:"many2many:room_admins"`
	AvatarURL   string
}

func (rm *RoomModel) GetIntID() int {
	return int(rm.ID)
}

func (rm *RoomModel) GetStringID() string {
	return ""
}

func (rm *RoomModel) GetName() string {
	return rm.Name
}

func (rm *RoomModel) GetAvatarURL() string {
	return rm.AvatarURL
}

func (rm *RoomModel) CreateLink() {
	id := strconv.Itoa(int(rm.ID))
	buf := &bytes.Buffer{}
	wc := base64.NewEncoder(base64.URLEncoding, buf)
	wc.Write([]byte(id + ":" + rm.Name))
	defer wc.Close()
	rm.Link = buf.String()
}

func (rm *RoomModel) AfterCreate(tx *gorm.DB) (err error) {
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
	GetAll() []*RoomModel
	FindMany(Map) []*RoomModel
	FindByID(id uint) *RoomModel
	FindByLink(string) *RoomModel
	Create(room *RoomModel) error
	Update(room *RoomModel) error
	Delete(id uint) error
	GetAdmins(*RoomModel) []*UserModel
	GetParticipants(*RoomModel) []*UserModel
	AddParticipant(*RoomModel, *UserModel) error
	RemoveParticipant(*RoomModel, *UserModel) error
}

type RoomGorm struct {
	*gorm.DB
}

func NewRoomGorm(db *gorm.DB) *RoomGorm {
	return &RoomGorm{db}
}

func (rg *RoomGorm) GetAll() []*RoomModel {
	var rm []*RoomModel
	rg.DB.Find(&rm)
	return rm
}

func (rg *RoomGorm) FindByLink(s string) *RoomModel {
	id := RoomIDFromLink(s)
	return rg.FindByID(uint(id))
}

func (rg *RoomGorm) FindByID(id uint) *RoomModel {
	return rg.byQuery(rg.DB.Where("id = ?", id))
}

func (rg *RoomGorm) FindMany(m Map) []*RoomModel {
	var rm []*RoomModel
	rg.DB.Where(m).Find(&rm)
	return rm
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

func (rg *RoomGorm) GetParticipants(r *RoomModel) []*UserModel {
	var users []*UserModel
	rg.DB.Model(r).Association("Users").Find(&users)
	return users
}

func (rg *RoomGorm) AddParticipant(room *RoomModel, u *UserModel) error {
	return rg.DB.Model(room).Association("Users").Append(u).Error
}

func (rg *RoomGorm) RemoveParticipant(room *RoomModel, u *UserModel) error {
	return rg.DB.Model(room).Association("Users").Delete(u).Error
}

func (rg *RoomGorm) GetAdmins(r *RoomModel) []*UserModel {
	var users []*UserModel
	rg.DB.Model(r).Association("Admins").Find(&users)
	return users
}

func (rg *RoomGorm) DestructiveReset() {
	rg.DropTableIfExists(&RoomModel{})
	rg.AutoMigrate()
}

func (rg *RoomGorm) AutoMigrate() {
	rg.DB.AutoMigrate(&RoomModel{})
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

package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Firstname string
	Lastname  string
	Nickname  string
	Email     string `gorm:"not null;unique_index"`
	Password  string
	AvatarURL string
	Rooms     []*Room    `gorm:"many2many:room_participants"`
	Messages  []*Message `gorm:"constraint:OnDelete:CASCADE;"`
}

func (um *User) GetIntID() int {
	return int(um.ID)
}

func (um *User) GetStringID() string {
	return ""
}

func (um *User) GetName() string {
	return um.Nickname
}

func (um *User) GetAvatarURL() string {
	return um.AvatarURL
}

// UserService is responsible for enabling communication between the model
// and the handlers
type UserService interface {
	Authenticate(email, password string) *User
	FindByID(id uint) *User
	FindByEmail(email string) *User
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
	GetRooms(*User) []*Room
}

// UserGorm reads and writes directly to the database
type UserGorm struct {
	*gorm.DB
}

func NewUserGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{db}
}

func (ug *UserGorm) FindByID(id uint) *User {
	return ug.byQuery(ug.DB.Where("id = ?", id))
}

func (ug *UserGorm) FindByEmail(email string) *User {
	return ug.byQuery(ug.DB.Where("email = ?", email))
}

func (ug *UserGorm) Create(user *User) error {
	return ug.DB.Create(user).Error
}

func (ug *UserGorm) Update(user *User) error {
	return ug.DB.Save(user).Error
}

func (ug *UserGorm) Delete(id uint) error {
	user := &User{}
	return ug.DB.Delete(user, id).Error
}
func (ug *UserGorm) GetRooms(u *User) []*Room {
	var rooms []*Room
	ug.DB.Model(u).Association("Rooms").Find(&rooms)
	return rooms
}

func (ug *UserGorm) byQuery(query *gorm.DB) *User {
	user := &User{}
	err := query.First(user).Error
	switch err {
	case nil:
		return user
	case gorm.ErrRecordNotFound:
		return nil
	default:
		panic(err)
	}
}

func (ug *UserGorm) DestructiveReset() {
	ug.DropTableIfExists(&User{})
	ug.AutoMigrate()
}

func (ug *UserGorm) AutoMigrate() {
	ug.DB.AutoMigrate(&User{})
}

func (ug *UserGorm) Authenticate(email, password string) *User {
	user := ug.FindByEmail(email)
	if user == nil || user.Password != password {
		return nil
	}
	return user
}

package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Firstname  string
	Lastname   string
	Username   string
	Email      string `gorm:"not null;unique_index"`
	Password   string
	AvatarURL  string
	Rooms      []*Room    `gorm:"many2many:room_participants;constraint:OnDelete:CASCADE;"`
	Messages   []*Message `gorm:"constraint:OnDelete:CASCADE;"`
	AdminRooms []*Room    `gorm:"many2many:room_admins;constraint:OnDelete:CASCADE;"`
}

func (u User) String() string {
	return (&u).GetFullname()
}

func (u *User) GetIntID() int {
	return int(u.ID)
}

func (u *User) GetFullname() string {
	return fmt.Sprintf("%s %s", u.Firstname, u.Lastname)
}

func (*User) GetStringID() string {
	return ""
}

func (u *User) GetName() string {
	return u.Username
}

func (u *User) GetAvatarURL() string {
	return u.AvatarURL
}

func (u *User) SetPasswordHash() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) VerifyPassword(p string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.SetPasswordHash()
	u.Email = strings.ToLower(u.Email)
	return nil
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
	user := ug.FindByEmail(strings.ToLower(email))
	if user == nil || user.VerifyPassword(password) != nil {
		return nil
	}
	return user
}

package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserModel struct {
	gorm.Model
	Firstname string
	Lastname  string
	Nickname  string
	Email     string `gorm:"not null;unique_index"`
	Password  string
	AvatarURL string
	Rooms     []*RoomModel `gorm:"many2many:room_participants"`
}

// UserService is responsible for enabling communication between the model
// and the handlers
type UserService interface {
	Authenticate(email, password string) *UserModel
	FindByID(id uint) *UserModel
	FindByEmail(email string) *UserModel
	Create(user *UserModel) error
	Update(user *UserModel) error
	Delete(id uint) error
	GetRooms(*UserModel) []*RoomModel
}

// UserGorm reads and writes directly to the database
type UserGorm struct {
	*gorm.DB
}

func NewUserGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{db}
}

func (ug *UserGorm) FindByID(id uint) *UserModel {
	return ug.byQuery(ug.DB.Where("id = ?", id))
}

func (ug *UserGorm) FindByEmail(email string) *UserModel {
	return ug.byQuery(ug.DB.Where("email = ?", email))
}

func (ug *UserGorm) Create(user *UserModel) error {
	return ug.DB.Create(user).Error
}

func (ug *UserGorm) Update(user *UserModel) error {
	return ug.DB.Save(user).Error
}

func (ug *UserGorm) Delete(id uint) error {
	user := ug.FindByID(id)
	return ug.DB.Delete(user).Error
}
func (ug *UserGorm) GetRooms(u *UserModel) []*RoomModel {
	var rooms []*RoomModel
	ug.DB.Model(u).Association("Rooms").Find(&rooms)
	return rooms
}

func (ug *UserGorm) byQuery(query *gorm.DB) *UserModel {
	user := &UserModel{}
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
	ug.DropTableIfExists(&UserModel{})
	ug.AutoMigrate()
}

func (ug *UserGorm) AutoMigrate() {
	ug.DB.AutoMigrate(&UserModel{})
}

func (ug *UserGorm) Authenticate(email, password string) *UserModel {
	user := ug.FindByEmail(email)
	if user == nil || user.Password != password {
		return nil
	}
	return user
}

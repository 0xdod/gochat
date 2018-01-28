package store

import (
	"context"

	"github.com/0xdod/gochat/gochat"

	"gorm.io/gorm"
)

type RoomGorm struct {
	*gorm.DB
}

func NewRoomStore(db *DB) *RoomGorm {
	return &RoomGorm{db.DB}
}

func (rs *RoomGorm) FindRoomByID(ctx context.Context, id int) (*gochat.Room, error) {
	room := &gochat.Room{}
	err := rs.DB.WithContext(ctx).First(room, id).Error

	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *RoomGorm) FindRooms(ctx context.Context, filter gochat.RoomFilter) ([]*gochat.Room, int, error) {
	var rooms []*gochat.Room
	err := rs.DB.WithContext(ctx).Find(&rooms).Error

	return rooms, 0, err
}

func (rs *RoomGorm) UpdateRoom(ctx context.Context, id int, upd gochat.RoomUpdate) (*gochat.Room, error) {
	return nil, nil
}

func (rs *RoomGorm) DeleteRoom(ctx context.Context, id int) error {
	return nil
}

func (rs *RoomGorm) CreateRoom(ctx context.Context, room *gochat.Room) error {
	err := rs.DB.WithContext(ctx).Create(room).Error

	if isDuplicateKeyError(err) {
		return gochat.ECONFLICT
	}
	return err
}

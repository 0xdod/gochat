package store

import (
	"context"

	"github.com/0xdod/gochat"

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
	err := rs.DB.WithContext(ctx).
		Preload("Participants").
		Preload("Admins").
		First(room, id).Error

	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *RoomGorm) FindRooms(ctx context.Context, filter gochat.RoomFilter) ([]*gochat.Room, int, error) {
	var rooms []*gochat.Room
	if v := filter.UserID; v != nil {
		return rs.findByParticipant(ctx, *v)
	}
	if v := filter.Invite; v != nil {
		return rs.findByInvite(ctx, *v)
	}
	err := rs.DB.WithContext(ctx).Find(&rooms).Error
	return rooms, 0, err
}

func (rs *RoomGorm) UpdateRoom(ctx context.Context, id int, upd gochat.RoomUpdate) (*gochat.Room, error) {
	return rs.updateRoom(ctx, id, upd)
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

func (rs *RoomGorm) findByInvite(ctx context.Context, invite string) ([]*gochat.Room, int, error) {
	var rooms []*gochat.Room
	err := rs.DB.WithContext(ctx).Where("invite = ?", invite).Find(&rooms).Error
	return rooms, 0, err
}

func (rs *RoomGorm) findByParticipant(ctx context.Context, id int) ([]*gochat.Room, int, error) {
	var rooms []*gochat.Room
	user := &gochat.User{}
	user.ID = id
	_ = rs.DB.WithContext(ctx).Model(user).
		Association("Rooms").Find(&rooms)
	return rooms, 0, nil
}

func (rs *RoomGorm) updateRoom(ctx context.Context, id int, upd gochat.RoomUpdate) (*gochat.Room, error) {
	room, err := rs.FindRoomByID(ctx, id)
	if err != nil {
		return room, err
	}
	if v := upd.Name; v != nil {
		room.Name = *v
	}
	if v := upd.Description; v != nil {
		room.Description = v
	}
	if v := upd.Link; v != nil {
		room.Link = *v
	}
	if v := upd.Icon; v != nil {
		room.Icon = v
	}
	if v := upd.Participants; v != nil {
		room.Participants = append(room.Participants, v...)
	}
	if v := upd.Admins; v != nil {
		room.Admins = append(room.Admins, v...)
	}
	return room, rs.DB.Save(room).Error
}

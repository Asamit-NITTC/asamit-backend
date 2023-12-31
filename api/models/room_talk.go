package models

import (
	"errors"
	"gorm.io/gorm"
)

type RoomTalk struct {
	RoomRoomID string `gorm:"default:null"`
	UserUID    string `gorm:"not null"`
	Comment    string `gorm:"not null"`
	ImageURL   string
	User       User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Room       Room `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	gorm.Model
}

type RoomTalkRepo struct {
	repo *gorm.DB
}

func InitializeRoomTaliRepo(db *gorm.DB) *RoomTalkRepo {
	return &RoomTalkRepo{repo: db}
}

type RoomTalkModel interface {
	InsertComment(rt RoomTalk) error
	GetAllTalk(roomId string) ([]RoomTalk, error)
	GetPersonalTalk(uid string) ([]RoomTalk, error)
}

func (r RoomTalkRepo) InsertComment(rt RoomTalk) error {
	err := r.repo.Create(&rt).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RoomTalkRepo) GetAllTalk(roomId string) ([]RoomTalk, error) {
	var roomTalkList []RoomTalk
	err := r.repo.Order("updated_at").Find(&roomTalkList, "room_room_id = ?", roomId).Error
	if err != nil {
		return roomTalkList, err
	}

	if len(roomTalkList) == 0 {
		return roomTalkList, errors.New("record not found.")
	}
	return roomTalkList, nil
}

func (r RoomTalkRepo) GetPersonalTalk(uid string) ([]RoomTalk, error) {
	var roomTalkList []RoomTalk
	err := r.repo.Where("user_uid = ?", uid).Where("room_room_id is null").Find(&roomTalkList).Error
	if err != nil {
		return roomTalkList, err
	}

	if len(roomTalkList) == 0 {
		return roomTalkList, errors.New("record not found.")
	}
	return roomTalkList, nil
}

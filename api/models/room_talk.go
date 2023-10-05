package models

import "gorm.io/gorm"

type RoomTalk struct {
	RoomRoomID string `gorm:"primaryKey"`
	UserUID    string `gorm:"primaryKey"`
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
}

func (r RoomTalkRepo) InsertComment(rt RoomTalk) error {
	err := r.repo.Create(&rt).Error
	if err != nil {
		return err
	}
	return nil
}

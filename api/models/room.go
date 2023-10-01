package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	RoomID     string
	TargetTime time.Time
	Decription string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type RoomRepo struct {
	db *gorm.DB
}

func InitializeRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{db: db}
}

type RoomModel interface {
	CreatRoom() error
}

func (r RoomRepo) CreateRoom() error {
	return nil
}

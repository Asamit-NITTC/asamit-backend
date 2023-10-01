package models

import "gorm.io/gorm"

type RoomUsersLink struct {
	RoomRoomID string `gorm:"unique;`
	UserUID    string `gorm:"unique;`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Room       Room   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	gorm.Model
}

type RoomUsersLinkRepo struct {
	db *gorm.DB
}

func InitializeRoomUsersLinkRepo(db *gorm.DB) *RoomUsersLinkRepo {
	return &RoomUsersLinkRepo{db: db}
}

type RoomUsersLinkModel interface {
	Insert() error
}

func (r RoomUsersLinkRepo) Insert() error {
	return nil
}

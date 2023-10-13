package models

import "gorm.io/gorm"

type RoomUsersLink struct {
	RoomRoomID string `gorm:"primaryKey"`
	UserUID    string `gorm:"primaryKey"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Room       Room   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type RoomUsersLinkRepo struct {
	repo *gorm.DB
}

func InitializeRoomUsersLinkRepo(db *gorm.DB) *RoomUsersLinkRepo {
	return &RoomUsersLinkRepo{repo: db}
}

type RoomUsersLinkModel interface {
	Insert(ru RoomUsersLink) error
	GetRoomIdIfAffiliated(uid string) (string, error)
	GetRoomBelongingUser(uid string) ([]RoomUsersLink, error)
}

func (r RoomUsersLinkRepo) Insert(ru RoomUsersLink) error {
	err := r.repo.Create(&ru).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RoomUsersLinkRepo) GetRoomIdIfAffiliated(uid string) (string, error) {
	var roomUserLinkInfo RoomUsersLink
	err := r.repo.Find(&roomUserLinkInfo, "user_uid = ?", uid).Error
	if err != nil {
		return "", err
	}
	return roomUserLinkInfo.RoomRoomID, nil
}

func (r RoomUsersLinkRepo) GetRoomBelongingUser(roomId string) ([]RoomUsersLink, error) {
	var roomUserLinkInfo []RoomUsersLink
	err := r.repo.Find(&roomUserLinkInfo, "room_room_id = ?", roomId).Error
	if err != nil {
		return roomUserLinkInfo, err
	}
	return roomUserLinkInfo, nil
}

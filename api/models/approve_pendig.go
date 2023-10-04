package models

import (
	"gorm.io/gorm"
)

type ApprovePendig struct {
	RoomRoomID string `gorm:"primaryKey"`
	UserUID    string `gorm:"primaryKey"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Room       Room   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ApprovePendigRepo struct {
	repo *gorm.DB
}

func InitializeApprovePendingRepo(db *gorm.DB) *ApprovePendigRepo {
	return &ApprovePendigRepo{repo: db}
}

type ApprovePendingModel interface {
	ReturnRoomIdIfRegisterd(uid string) (string, error)
}

func (a ApprovePendigRepo) ReturnRoomIdIfRegisterd(uid string) (string, error) {
	var approvePending ApprovePendig
	err := a.repo.Find(&approvePending, "user_uid = ?", uid).Error
	if err != nil {
		return "", err
	}
	return approvePending.RoomRoomID, nil
}

package models

import (
	"gorm.io/gorm"
)

type ApprovePending struct {
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
	CheckExists(uid string) (bool, error)
	DeletePendingRecord(uid string) error
	GetRoomId(uid string) (string, error)
	InsertApprovePendingUserList(ap []ApprovePending) error
}

func (a ApprovePendigRepo) ReturnRoomIdIfRegisterd(uid string) (string, error) {
	var approvePending ApprovePending
	err := a.repo.Find(&approvePending, "user_uid = ?", uid).Error
	if err != nil {
		return "", err
	}
	return approvePending.RoomRoomID, nil
}

func (a ApprovePendigRepo) CheckExists(uid string) (bool, error) {
	var approvePending ApprovePending
	err := a.repo.Find(&approvePending, "user_uid = ?", uid).Error
	if err != nil {
		return false, err
	}

	if approvePending.RoomRoomID == "" {
		return false, nil
	}
	return true, nil
}

func (a ApprovePendigRepo) DeletePendingRecord(uid string) error {
	err := a.repo.Delete(&ApprovePending{}, "user_uid = ?", uid).Error
	if err != nil {
		return err
	}
	return nil
}

func (a ApprovePendigRepo) GetRoomId(uid string) (string, error) {
	var approvePendingInfo ApprovePending
	err := a.repo.Find(&approvePendingInfo, "user_uid = ?", uid).Error
	if err != nil {
		return "", err
	}
	return approvePendingInfo.RoomRoomID, nil
}

func (a ApprovePendigRepo) InsertApprovePendingUserList(ap []ApprovePending) error {
	err := a.repo.Create(&ap).Error
	if err != nil {
		return err
	}
	return nil
}

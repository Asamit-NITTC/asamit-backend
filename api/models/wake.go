package models

import (
	"time"

	"gorm.io/gorm"
)

type Wake struct {
	UserUID    string `json:"uid" gorm:"size:256"`
	RoomRoomID string
	WakeUpTime string `json:"WakeUpTime"`
	Comment    string `json:"comment"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Room       Room   `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	gorm.Model
}

type WakeRepo struct {
	repo *gorm.DB
}

func InitializeWakeRepo(repo *gorm.DB) *WakeRepo {
	return &WakeRepo{repo: repo}
}

type WakeModel interface {
	Report(wa Wake) error
	GetAllReport(uid string) ([]Wake, error)
}

func (w WakeRepo) Report(wa Wake) error {
	_, err := time.Parse(time.RFC3339, wa.WakeUpTime)
	if err != nil {
		return err
	}
	err = w.repo.Save(&wa).Where("uid = ?", wa.UserUID).Error
	if err != nil {
		return err
	}
	return nil
}

func (w WakeRepo) GetAllReport(uid string) ([]Wake, error) {
	var allWakeUpReport []Wake
	err := w.repo.Find(&allWakeUpReport, "user_uid = ?", uid).Order("updated_at").Error
	if err != nil {
		return allWakeUpReport, err
	}
	return allWakeUpReport, nil
}

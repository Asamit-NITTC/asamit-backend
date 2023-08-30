package models

import (
	"time"

	"gorm.io/gorm"
)

type Wake struct {
	ReportID   int    `gorm:"primaryKey;autoIncrement"`
	UserUID    string `json:"uid" gorm:"size:256"`
	WakeUpTime string `json:"WakeUpTime"`
	Comment    string `json:"comment"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type WakeRepo struct {
	repo *gorm.DB
}

func InitalizeWakeRepo(repo *gorm.DB) *WakeRepo {
	return &WakeRepo{repo: repo}
}

type WakeModel interface {
	Report(wa Wake) error
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

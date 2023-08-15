package models

import (
	"time"

	"github.com/Asamit-NITTC/asamit-backend-test/db"
)

type Wake struct {
	ReportID   int    `gorm:"primaryKey;autoIncrement"`
	UserUID    string `json:"uid" gorm:"size:256"`
	WakeUpTime string `json:"WakeUpTime"`
	Comment    string `json:"comment"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type WakeModel struct{}

func (w WakeModel) Report(wa Wake) error {
	_, err := time.Parse(time.RFC3339, wa.WakeUpTime)
	if err != nil {
		return err
	}
	err = db.DB.Save(&wa).Where("uid = ?", wa.UserUID).Error
	if err != nil {
		return err
	}
	return nil
}

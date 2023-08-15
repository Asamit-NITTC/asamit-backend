package models

import (
	"time"

	"github.com/Asamit-NITTC/asamit-backend-test/db"
)

type WakeUpTime struct {
	UserUID    string `json:"uid" gorm:"primaryKey;size:256"`
	TargetTime string `json:"targetTime"`
	Updated    int64  `json:"updated" gorm:"autoUpdateTime:nano"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type WakeModel struct{}

func (w WakeModel) Set(wt WakeUpTime) error {
	_, err := time.Parse(time.RFC3339, wt.TargetTime)
	if err != nil {
		return err
	}
	err = db.DB.Save(&wt).Where("uid = ?", wt.UserUID).Error
	if err != nil {
		return err
	}
	return nil
}

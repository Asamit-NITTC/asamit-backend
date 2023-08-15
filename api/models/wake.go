package models

import (
	"time"

	"github.com/Asamit-NITTC/asamit-backend-test/db"
)

type WakeUpTime struct {
	//UID        []Users `json:"uid" gorm:"foreignKey:UID"`
	UID        []User `json:"uid"`
	TargetTime string `json:"targetTime" `
	Updated    int64  `json:"updated" gorm:"autoUpdateTime:nano"`
}

type WakeModel struct{}

func (w WakeModel) Set(wt WakeUpTime) error {
	_, err := time.Parse(time.RFC3339, wt.TargetTime)
	if err != nil {
		return err
	}
	err = db.DB.Save(&wt).Where("uid = ?", wt.UID).Error
	if err != nil {
		return err
	}
	return nil
}

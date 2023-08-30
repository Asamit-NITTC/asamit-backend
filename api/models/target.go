package models

import (
	"time"

	"gorm.io/gorm"
)

type TargetTime struct {
	UserUID    string `json:"uid" gorm:"primaryKey;size:256"`
	TargetTime string `json:"targetTime"`
	Updated    int64  `json:"updated" gorm:"autoUpdateTime:nano"`
	User       User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type TargetTimeRepo struct {
	db *gorm.DB
}

func InitializeTargetRepo(db *gorm.DB) *TargetTimeRepo {
	return &TargetTimeRepo{db: db}
}

type TargetTimeModel interface {
	Set(wt TargetTime) error
}

func (t TargetTimeRepo) Set(wt TargetTime) error {
	_, err := time.Parse(time.RFC3339, wt.TargetTime)
	if err != nil {
		return err
	}
	err = t.db.Save(&wt).Where("uid = ?", wt.UserUID).Error
	if err != nil {
		return err
	}
	return nil
}

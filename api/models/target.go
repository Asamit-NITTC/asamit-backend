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
	repo *gorm.DB
}

func InitializeTargetRepo(db *gorm.DB) *TargetTimeRepo {
	return &TargetTimeRepo{repo: db}
}

type TargetTimeModel interface {
	Set(wt TargetTime) error
	Get(uid string) (string, error)
}

func (t TargetTimeRepo) Set(wt TargetTime) error {
	_, err := time.Parse(time.RFC3339, wt.TargetTime)
	if err != nil {
		return err
	}
	err = t.repo.Save(&wt).Where("uid = ?", wt.UserUID).Error
	if err != nil {
		return err
	}
	return nil
}

func (t TargetTimeRepo) Get(uid string) (string, error) {
	var targetTimeInfo TargetTime
	err := t.repo.First(&targetTimeInfo).Error
	if err != nil {
		return "", err
	}
	return targetTimeInfo.TargetTime, nil
}

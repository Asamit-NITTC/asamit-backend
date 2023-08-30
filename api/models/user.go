package models

import (
	"gorm.io/gorm"
)

type User struct {
	UID      string `json:"uid" gorm:"primaryKey;size:256"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Point    int    `json:"point"`
	Duration int    `json:"duration"`
}

type UserRepo struct {
	repo *gorm.DB
}

func InitalizeUserRepo(repo *gorm.DB) *UserRepo {
	return &UserRepo{repo: repo}
}

type UserModel interface {
	GetUserInfo(uid string) (User, error)
	SetUserInfo(us *User) error
}

func (u UserRepo) GetUserInfo(uid string) (User, error) {
	var userInfo User
	err := u.repo.First(&userInfo, "uid = ?", uid).Error
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (u UserRepo) SetUserInfo(us *User) error {
	err := u.repo.Save(&us).Where("uid = ?", us.UID).Error
	if err != nil {
		return err
	}
	return nil
}

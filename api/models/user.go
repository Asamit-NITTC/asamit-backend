package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	UID      string `json:"uid" gorm:"primaryKey;not null;size:256"`
	Sub      string `json:"sub" gorm:"primaryKey;not null;size:500"`
	Name     string `json:"name"`
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
	SignUpUserInfo(us *User) error
	ChangeUserInfo(us User) error
	CheckExistsUser(uid string) (string, error)
}

func (u UserRepo) GetUserInfo(uid string) (User, error) {
	var userInfo User
	err := u.repo.First(&userInfo, "uid = ?", uid).Error
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (u UserRepo) SignUpUserInfo(us *User) error {
	err := u.repo.Save(&us).Where("uid = ?", us.UID).Error
	if err != nil {
		return err
	}
	return nil
}

func (u UserRepo) ChangeUserInfo(us User) error {
	err := u.repo.Model(&User{}).Where("uid = ?", us.UID).UpdateColumns(User{Name: us.Name, Point: us.Point, Duration: us.Duration}).Error
	if err != nil {
		return err
	}
	return nil
}

func (u UserRepo) CheckExistsUser(uid string) (string, error) {
	var userInfo User
	err := u.repo.First(&userInfo, "uid = ?", uid).Error
	if err != nil {
		return "", err
	}

	if uid == "" {
		return "", errors.New("Subが空です")
	}
	return userInfo.Sub, nil
}

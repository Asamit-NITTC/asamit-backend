package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UID      string `json:"uid" gorm:"primaryKey;not null;size:256"`
  Sub      string `json:"sub" gorm:"unique;not null;size: 600"`
	Name     string `json:"name"`
	Point    int    `json:"point"`
	Duration int    `json:"duration"`
}

type UserRepo struct {
	repo *gorm.DB
}

func InitializeUserRepo(repo *gorm.DB) *UserRepo {
	return &UserRepo{repo: repo}
}

type UserModel interface {
	GetUserInfo(uid string) (User, error)
	SignUpUserInfo(us *User) error
	ChangeUserInfo(us User) error
	CheckExistsUserWithUID(uid string) (string, error)
	CheckExistsUserWithSub(sub string) (bool, error)
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
	us.UID = uuid.NewString()
	err := u.repo.Create(us).Error
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

func (u UserRepo) CheckExistsUserWithUID(uid string) (string, error) {
	var userInfo User
	err := u.repo.First(&userInfo, "uid = ?", uid).Error
	if err != nil {
		return "", err
	}

	return userInfo.Sub, nil
}

func (u UserRepo) CheckExistsUserWithSub(sub string) (bool, error) {
	var userInfo User
	err := u.repo.Find(&userInfo, "sub = ?", sub).Error
	if err != nil {
		return false, err
	}

	//あくまでも判定はController層に委ねる
	if userInfo.Sub == "" {
		return false, nil
	}

	return true, nil
}

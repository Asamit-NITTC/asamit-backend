package models

import "github.com/Asamit-NITTC/asamit-backend-test/db"

type User struct {
	UID      string `json:"uid" gorm:"primaryKey;size:256"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Point    int    `json:"point"`
	Duration int    `json:"duration"`
}

type UsersModel struct{}

func (u UsersModel) GetUserInfo(uid string) (User, error) {
	var userInfo User
	err := db.DB.First(&userInfo, "uid = ?", uid).Error
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (u UsersModel) SetUserInfo(us *User) error {
	err := db.DB.Save(&us).Where("uid = ?", us.UID).Error
	if err != nil {
		return err
	}
	return nil
}

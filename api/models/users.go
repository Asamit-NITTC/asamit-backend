package models

import "github.com/Asamit-NITTC/asamit-backend-test/db"

type Users struct {
	UID      string `json:"uid" gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Icon     string `json:"icon"`
	Point    int    `json:"point"`
	Duration int    `json:"duration"`
}

type UsersModel struct{}

func (u UsersModel) GetUserInfo(uid string) (Users, error) {
	var userInfo Users
	err := db.DB.First(&userInfo, "uid = ?", uid).Error
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

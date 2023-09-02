package models

import (
	"gorm.io/gorm"
)

type AuthRepo struct {
	repo *gorm.DB
}

func InitializeAuthRepo(repo *gorm.DB) *AuthRepo {
	return &AuthRepo{repo: repo}
}

type AuthModel interface {
	CheckSubIsValid(sub string) (bool, error)
}

func (u AuthRepo) CheckSubIsValid(sub string) (bool, error) {
	var user User
	result := u.repo.Where("sub = ?", sub).Where(&user)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected != 1 {
		return false, nil
	}

	return true, nil
}

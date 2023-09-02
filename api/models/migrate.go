package models

import (
	"os"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&User{}, &TargetTime{}, &Wake{})
}

func InsertDummyData(db *gorm.DB) {
	var users = []User{
		{UID: "33u@2", Sub: os.Getenv("TEST_SUB"), Name: "GoRuGoo", Icon: "https://upload.wikimedia.org/wikipedia/en/thumb/5/5f/Original_Doge_meme.jpg/300px-Original_Doge_meme.jpg", Point: 32, Duration: 5},
	}
	db.Save(&users)
}

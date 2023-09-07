package models

import (
	"os"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	db.Migrator().DropTable()
	db.AutoMigrate(&User{}, &TargetTime{}, &Wake{})
}

func InsertDummyData(db *gorm.DB) {
	var users = []User{
		{UID: "33u@2", Sub: os.Getenv("TEST_SUB"), Name: "GoRuGoo", Point: 32, Duration: 5},
	}
	db.Save(&users)
}

package models

import "github.com/Asamit-NITTC/asamit-backend-test/db"

func MigrateDB() {
	db.DB.AutoMigrate(&Users{})
	db.DB.AutoMigrate(&WakeUpTime{})
}

func InsertDummyData() {
	var users = []Users{
		{UID: "33u@2", Name: "GoRuGoo", Icon: "https://upload.wikimedia.org/wikipedia/en/thumb/5/5f/Original_Doge_meme.jpg/300px-Original_Doge_meme.jpg", Point: 32, Duration: 5},
	}
	db.DB.Save(&users)
}

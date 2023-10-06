package models

import (
	"os"
	"time"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	db.Migrator().DropTable(&User{}, &TargetTime{}, &Wake{}, &Room{}, &RoomUsersLink{}, &ApprovePendig{})
	db.AutoMigrate(&User{}, &TargetTime{}, &Wake{}, &Room{}, &RoomUsersLink{}, &ApprovePendig{})
}

func InsertDummyData(db *gorm.DB) {
	var users = []User{
		{UID: "33u@2", Sub: os.Getenv("TEST_SUB"), Name: "GoRuGoo", Point: 32, Duration: 5, InvitationStatus: false, AffiliationStatus: true},
		{UID: "xyz123", Sub: "abc123", Name: "Alice", Point: 45, Duration: 8, InvitationStatus: true, AffiliationStatus: false},
		{UID: "789abc", Sub: "def456", Name: "Bob", Point: 27, Duration: 3, InvitationStatus: true, AffiliationStatus: false},
		{UID: "ghi789", Sub: "jkl012", Name: "Charlie", Point: 19, Duration: 6, InvitationStatus: true, AffiliationStatus: false},
		{UID: "321jkl", Sub: "mno345", Name: "David", Point: 55, Duration: 9, InvitationStatus: true, AffiliationStatus: false},
		{UID: "456pqr", Sub: "stu789", Name: "Eve", Point: 12, Duration: 2, InvitationStatus: true, AffiliationStatus: false},
		{UID: "lmn012", Sub: "vwx345", Name: "Frank", Point: 60, Duration: 7, InvitationStatus: true, AffiliationStatus: false},
		{UID: "def345", Sub: "yza678", Name: "Grace", Point: 36, Duration: 4, InvitationStatus: false, AffiliationStatus: false},
		{UID: "hij678", Sub: "bcd901", Name: "Hank", Point: 25, Duration: 10, InvitationStatus: false, AffiliationStatus: false},
		{UID: "123bcd", Sub: "efg234", Name: "Ivy", Point: 42, Duration: 3, InvitationStatus: true, AffiliationStatus: false},
	}

	rfc3339FormattedCurrentTime := time.Now().Format(time.RFC3339)

	var targetTime = []TargetTime{
		{UserUID: "33u@2", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "xyz123", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "789abc", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "ghi789", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "123bcd", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "lmn012", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "def345", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "hij678", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "456pqr", TargetTime: rfc3339FormattedCurrentTime},
		{UserUID: "321jkl", TargetTime: rfc3339FormattedCurrentTime},
	}

	var wakeData = []Wake{
		{UserUID: "33u@2", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Good morning!"},
		{UserUID: "xyz123", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Rise and shine!"},
		{UserUID: "789abc", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Waking up early today."},
		{UserUID: "ghi789", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Another day begins."},
		{UserUID: "123bcd", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Ready to start the day."},
		{UserUID: "lmn012", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Morning routine."},
		{UserUID: "def345", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Feeling refreshed!"},
		{UserUID: "hij678", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Early bird."},
		{UserUID: "456pqr", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "A new dawn."},
		{UserUID: "321jkl", WakeUpTime: rfc3339FormattedCurrentTime, Comment: "Good morning world!"},
	}

	convetedTime, _ := time.Parse(time.RFC3339, rfc3339FormattedCurrentTime)
	var room = []Room{
		{RoomID: "ohayou", WakeUpTime: convetedTime, Decription: "test"},
	}

	var roomUserLink = []RoomUsersLink{
		{RoomRoomID: "ohayou", UserUID: "33u@2"},
	}

	var approvePending = []ApprovePendig{
		{RoomRoomID: "ohayou", UserUID: "123bcd"},
	}

	db.Save(&users)
	db.Save(&targetTime)
	db.Save(&wakeData)
	db.Save(&room)
	db.Save(&roomUserLink)
	db.Save(&approvePending)
}

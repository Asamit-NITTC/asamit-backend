package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	RoomID     string
	WakeUpTime time.Time
	Decription string
	gorm.Model
}

type RoomRepo struct {
	db *gorm.DB
}

func InitializeRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{db: db}
}

type RoomModel interface {
	CreatRoom(ro Room) (Room, error)
}

func (r RoomRepo) CreateRoom(ro Room) (Room, error) {
	//返却・書き込み用構造体
	roomInfoResult := ro

	formattedTime := ro.WakeUpTime.Format(time.RFC3339)
	formattedRFC3399TypeTime, err := time.Parse(time.RFC3339, formattedTime)
	if err != nil {
		return roomInfoResult, err
	}
	roomInfoResult.WakeUpTime = formattedRFC3399TypeTime

	roomId := uuid.NewString()
	//DBに書き込むためにUUIDをここで生成してRoomIDとする
	roomInfoResult.RoomID = roomId

	r.db.Create(&roomInfoResult)

	//後の中間テーブルに書き込むためにRoomIDを含む構造体を返す
	return roomInfoResult, nil
}

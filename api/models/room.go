package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	RoomID      string    `json:"roomID" gorm:"primaryKey;size:256"`
	WakeUpTime  time.Time `gorm:"not null"`
	Description string
	Mission     string `json:"mission" gorm:"not null"`
}

type RoomRepo struct {
	repo *gorm.DB
}

func InitializeRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{repo: db}
}

type RoomModel interface {
	CreateRoom(ro Room) (Room, error)
	GetRoomDetailInfo(roomID string) (Room, error)
	ChangeMission(ro Room) error
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

	r.repo.Create(&roomInfoResult)

	//後の中間テーブルに書き込むためにRoomIDを含む構造体を返す
	return roomInfoResult, nil
}

func (r RoomRepo) GetRoomDetailInfo(roomID string) (Room, error) {
	var roomInfo Room
	err := r.repo.First(&roomInfo, "room_id = ?", roomID).Error
	if err != nil {
		return roomInfo, err
	}

	return roomInfo, nil
}

func (r RoomRepo) ChangeMission(ro Room) error {
	if ro.RoomID == "" || ro.Mission == "" {
		return errors.New("roomId or mission is empty.")
	}
	err := r.repo.Model(&Room{}).Where("room_id = ?", ro.RoomID).Update("mission", ro.Mission).Error
	if err != nil {
		return err
	}
	return nil
}

package controllers

import (
	"net/http"
	"time"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type SummitController struct {
	roomModel          models.RoomModel
	userModel          models.UserModel
	roomUsersLinkModel models.RoomUsersLinkModel
}

func InitailizeRoomController(r models.RoomModel, u models.UserModel, ru models.RoomUsersLinkModel) *SummitController {
	return &SummitController{roomModel: r, userModel: u, roomUsersLinkModel: ru}
}

type createRoomRequestBody struct {
	MemberUID   []string  `json:"memberUID"`
	WakeUpTime  time.Time `json:"wakeUpTime"`
	Description string    `json:"description"`
}

func (s SummitController) Create(c *gin.Context) {
	var requestBody createRoomRequestBody
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, "Can't convert to json."})
		return
	}

	//DBに入れる各種情報の書き込み
	var roomInfo models.Room
	roomInfo.WakeUpTime = requestBody.WakeUpTime
	roomInfo.Decription = requestBody.Description

	createdRoomInfo, err := s.roomModel.CreateRoom(roomInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error()})
		return
	}

	for _, uid := range requestBody.MemberUID {
		existUID, err := s.userModel.CheckExistsUserWithUIDReturnBool(uid)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error()})
			return
		}

		if !existUID {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusNotFound, "There are unregistered users."})
			return
		}

		//中間テーブル書き込み用
		var roomUsersLink models.RoomUsersLink
		roomUsersLink.RoomRoomID = createdRoomInfo.RoomID
		roomUsersLink.UserUID = uid

		err = s.roomUsersLinkModel.Insert(roomUsersLink)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, requestBody)
}

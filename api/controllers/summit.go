package controllers

import (
	"net/http"
	"time"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type SummitController struct {
	roomModel models.RoomModel
	userModel models.UserModel
}

func InitailizeSummitController(r models.RoomModel, u models.UserModel) *SummitController {
	return &SummitController{roomModel: r, userModel: u}
}

type createRoomRequestBody struct {
	MemberUID   []string
	WakeUpTime  time.Time
	Description string
}

func (s SummitController) CreateRoom(c *gin.Context) {
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

	createdRoomInfo, err := s.roomModel.CreatRoom(roomInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error()})
		return
	}
}

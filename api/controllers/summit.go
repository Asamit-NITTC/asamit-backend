package controllers

import (
	"net/http"
	"time"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type SummitController struct {
	roomModel           models.RoomModel
	userModel           models.UserModel
	roomUsersLinkModel  models.RoomUsersLinkModel
	approvePendingModel models.ApprovePendingModel
}

func InitailizeRoomController(r models.RoomModel, u models.UserModel, ru models.RoomUsersLinkModel, a models.ApprovePendingModel) *SummitController {
	return &SummitController{roomModel: r, userModel: u, roomUsersLinkModel: ru, approvePendingModel: a}
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
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "Can't convert to json."})
		return
	}

	//DBに入れる各種情報の書き込み
	var roomInfo models.Room
	roomInfo.WakeUpTime = requestBody.WakeUpTime
	roomInfo.Decription = requestBody.Description

	createdRoomInfo, err := s.roomModel.CreateRoom(roomInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Can't get RoomInfo."})
		return
	}

	for _, uid := range requestBody.MemberUID {
		existUID, err := s.userModel.CheckExistsUserWithUIDReturnBool(uid)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Can't get UID."})
			return
		}

		if !existUID {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusNotFound, "There are unregistered users.", "There are unregistered users."})
			return
		}

		//中間テーブル書き込み用
		var roomUsersLink models.RoomUsersLink
		roomUsersLink.RoomRoomID = createdRoomInfo.RoomID
		roomUsersLink.UserUID = uid

		err = s.roomUsersLinkModel.Insert(roomUsersLink)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
			return
		}
	}

	c.JSON(http.StatusOK, requestBody)
}

func (s SummitController) CheckAffiliateAndInventionStatus(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect query parameter."})
	}
	//本来はここに認証があってもいいが、現在の仕様はAuthorizationMiddlewareに一任している

	invitationStatus, err := s.userModel.CheckInvitationStatus(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
		return
	}

	if invitationStatus {
		roomID, err := s.approvePendingModel.ReturnRoomIdIfRegisterd(uid)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"roomID": roomID, "status": "Approval pending"})
	}

	affiliationStatus, err := s.userModel.CheckAffliationStatus(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
		return
	}

	if affiliationStatus {
		roomID, err := s.roomUsersLinkModel.GetRoomIdIfAffiliated(uid)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"roomID": roomID, "status": "Already belonging"})
	}

	c.JSON(http.StatusOK, gin.H{"roomID": "", "status": "Not using summit mode"})
}

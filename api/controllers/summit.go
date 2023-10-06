package controllers

import (
	"net/http"
	"time"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/Asamit-NITTC/asamit-backend-test/webstorage"
	"github.com/gin-gonic/gin"
)

type SummitController struct {
	roomModel            models.RoomModel
	userModel            models.UserModel
	roomUsersLinkModel   models.RoomUsersLinkModel
	approvePendingModel  models.ApprovePendingModel
	roomTalkModel        models.RoomTalkModel
	cloudStorageWebModel webstorage.CloudStorageOriginalWebModel
}

func InitailizeRoomController(r models.RoomModel, u models.UserModel, ru models.RoomUsersLinkModel, a models.ApprovePendingModel, rt models.RoomTalkModel, c webstorage.CloudStorageOriginalWebModel) *SummitController {
	return &SummitController{roomModel: r, userModel: u, roomUsersLinkModel: ru, approvePendingModel: a, roomTalkModel: rt, cloudStorageWebModel: c}
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

		if roomID == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "DB get error."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"roomID": roomID, "status": "Approval pending"})
		return
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

		if roomID == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "DB get error."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"roomID": roomID, "status": "Already belonging"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roomID": "", "status": "Not using summit mode"})
}

// 時間があったらテーブル結合を用いて実装したい
func (s SummitController) GetRoomDetailInfo(c *gin.Context) {
	roomID := c.Query("roomID")
	if roomID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect query parameter."})
		return
	}

	roomDetailInfo, err := s.roomModel.GetRoomDetailInfo(roomID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "DB get error."})
		return
	}

	c.JSON(http.StatusOK, roomDetailInfo)
}

func (s SummitController) RecordTalk(c *gin.Context) {
	var requestBody models.RoomTalk
	requestBody.RoomRoomID = c.PostForm("room_room_id")
	requestBody.UserUID = c.PostForm("user_uid")
	requestBody.Comment = c.PostForm("comment")
	//わかりやすい変数名に変更
	forWritingCommentObject := requestBody

	morningActivityImageFile, _ := c.FormFile("image")
	//下でバリデーションしているためあえてerrを受け取らない
	//ファイルサイズが0ならそもそもファイル関連の処理が走らないから安全

	if morningActivityImageFile.Size != 0 {
		morningActivityImage, err := morningActivityImageFile.Open()
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Image processing error."})
			return
		}
		defer morningActivityImage.Close()
		objectName, err := s.cloudStorageWebModel.Write(requestBody.RoomRoomID, morningActivityImage)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "Cloud Storage error."})
			return
		}
		forWritingCommentObject.ImageURL = objectName
	}
	err := s.roomTalkModel.InsertComment(forWritingCommentObject)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}
	c.JSON(http.StatusOK, forWritingCommentObject)
}

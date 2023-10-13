package controllers

import (
	"net/http"
	"time"

	"fmt"
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
	HostUID     string    `json:"hostUID"`
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
	roomInfo.Description = requestBody.Description

	//デフォルトのミッションを書き込む
	roomInfo.Mission = "今日の天気予報の写真を撮影せよ！"

	//ルーム作成(ユーザー関連操作なし)
	createdRoomInfo, err := s.roomModel.CreateRoom(roomInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Can't get RoomInfo."})
		return
	}

	//ホストは条件なしにルームに入れる
	//存在確認・ステータス変更
	hostUID := requestBody.HostUID
	existUID, err := s.userModel.CheckExistsUserWithUIDReturnBool(hostUID)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Can't get UID."})
		return
	}
	if !existUID {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found."})
	}
	err = s.userModel.ChangeInvitationAndAffiliationStatus(hostUID, false, true)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}

	//ホストを所属させる
	var roomUserLink models.RoomUsersLink
	roomUserLink.RoomRoomID = createdRoomInfo.RoomID
	roomUserLink.UserUID = hostUID

	err = s.roomUsersLinkModel.Insert(roomUserLink)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}

	var approvePendingUserTmp models.ApprovePending
	var approvePendingUserList []models.ApprovePending

	for _, uid := range requestBody.MemberUID {
		//ユーザー登録されているか確認する
		existUID, err := s.userModel.CheckExistsUserWithUIDReturnBool(uid)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Can't get UID."})
			return
		}
		if !existUID {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusNotFound, "There are unregistered users.", "There are unregistered users."})
			return
		}

		//既にどこかのルームに所属していないか確認する
		isAffiliated, err := s.approvePendingModel.CheckExists(uid)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
			return
		}

		if isAffiliated {
			fmt.Println("koko")
			fmt.Println(uid)
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Already belongs to the room."})
			return
		}

		//for後に承認待ちテーブルに一気に書き込む用の配列
		approvePendingUserTmp.RoomRoomID = createdRoomInfo.RoomID
		approvePendingUserTmp.UserUID = uid
		approvePendingUserList = append(approvePendingUserList, approvePendingUserTmp)

		err = s.userModel.ChangeInvitationAndAffiliationStatus(uid, true, false)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
			return
		}
	}

	err = s.approvePendingModel.InsertApprovePendingUserList(approvePendingUserList)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}

	c.JSON(http.StatusOK, requestBody)
}

func (s SummitController) CheckAffiliateAndInvitationStatus(c *gin.Context) {
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

func (s SummitController) Approve(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request."})
		return
	}

	isWatingAffiliation, err := s.approvePendingModel.CheckExists(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "DB get error."})
		return
	}

	if !isWatingAffiliation {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Can't find pending list."})
		return
	}

	roomId, err := s.approvePendingModel.GetRoomId(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
		return
	}

	err = s.userModel.ChangeInvitationAndAffiliationStatus(uid, false, true)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Infomation delete error."})
		return
	}

	err = s.approvePendingModel.DeletePendingRecord(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Infomation delete error."})
		return
	}

	var roomUsersLink models.RoomUsersLink
	roomUsersLink.RoomRoomID = roomId
	roomUsersLink.UserUID = uid
	err = s.roomUsersLinkModel.Insert(roomUsersLink)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
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

func (s SummitController) GetTalk(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UID and roomId are empty."})
	}
	roomId := c.Query("room-id")
	//個人・サミット両方の機能を実装するため片方が空だからすぐにエラーを返すわけじゃない

	if roomId != "" {
		//サミットの振り返り取得
		affiliateUID, err := s.roomUsersLinkModel.GetRoomIdIfAffiliated(uid)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
			return
		}

		//指定されたRoomIdと所属しているRoomIdが違ったらエラーを返す
		if affiliateUID != roomId {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Affiliation and request are different."})
			return
		}

		roomTalkList, err := s.roomTalkModel.GetAllTalk(roomId)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
			return
		}

		c.JSON(http.StatusOK, roomTalkList)
	}

	//個人の記録取得
	personalTalkList, err := s.roomTalkModel.GetPersonalTalk(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
		return
	}

	c.JSON(http.StatusOK, personalTalkList)
}

func (s SummitController) GetRoomUserLists(c *gin.Context) {
	roomId := c.Query("room-id")
	if roomId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "UID is empty."})
		return
	}
	roomBelongingUserList, err := s.roomUsersLinkModel.GetRoomBelongingUser(roomId)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
		return
	}

	c.JSON(http.StatusOK, roomBelongingUserList)
}

func (s SummitController) ChangeMission(c *gin.Context) {
	var changeTargetMissionInfo models.Room
	err := s.roomModel.ChangeMission(changeTargetMissionInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}
	c.JSON(http.StatusOK, changeTargetMissionInfo)
}

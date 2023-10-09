package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type WakeController struct {
	wakeModel          models.WakeModel
	userModel          models.UserModel
	roomusersLinkModel models.RoomUsersLinkModel
}

func InitializeWakeController(w models.WakeModel, u models.UserModel, r models.RoomUsersLinkModel) *WakeController {
	return &WakeController{wakeModel: w, userModel: u, roomusersLinkModel: r}
}

func (w WakeController) Report(c *gin.Context) {
	var wakeUpInfo models.Wake
	err := c.ShouldBindJSON(&wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "Can't convert to json."})
		return
	}

	roomId, err := w.roomusersLinkModel.GetRoomIdIfAffiliated(wakeUpInfo.UserUID)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "DB get error."})
		return
	}

	//空でも大丈夫
	wakeUpInfo.RoomRoomID = roomId

	err = w.wakeModel.Report(wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}
	c.JSON(http.StatusOK, wakeUpInfo)
	return
}

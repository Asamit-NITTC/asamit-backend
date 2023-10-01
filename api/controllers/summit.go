package controllers

import (
	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SummitController struct {
	roomModel models.RoomModel
	userModel models.UserModel
}

func InitailizeSummitController(r models.RoomModel, u models.UserModel) *SummitController {
	return &SummitController{roomModel: r, userModel: u}
}

type createRoomRequestBody struct {
	MemberUID  []string
	WakeUpTime string
}

func (s SummitController) CreateRoom(c *gin.Context) {
	var requestBody createRoomRequestBody
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, "Can't convert to json."})
		return
	}
}

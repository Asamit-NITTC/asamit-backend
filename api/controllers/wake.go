package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type WakeController struct {
	wakeModel models.WakeModel
	userModel models.UserModel
}

func InitializeWakeController(w models.WakeModel, u models.UserModel) *WakeController {
	return &WakeController{wakeModel: w, userModel: u}
}

func (w WakeController) Report(c *gin.Context) {
	var wakeUpInfo models.Wake
	err := c.ShouldBindJSON(&wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "Can't convert to json."})
		return
	}

	err = w.wakeModel.Report(wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}
	c.JSON(http.StatusOK, wakeUpInfo)
	return
}

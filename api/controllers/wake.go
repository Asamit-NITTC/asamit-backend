package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type WakeController struct {
	wakeModel models.WakeModel
}

func InitializeWakeController(w models.WakeModel) *WakeController {
	return &WakeController{wakeModel: w}
}

func (w WakeController) Report(c *gin.Context) {
	var wakeUpInfo models.Wake
	err := c.ShouldBindJSON(&wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	err = w.wakeModel.Report(wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, wakeUpInfo)
	return
}

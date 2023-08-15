package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type WakeController struct{}

var wakeModel = new(models.WakeModel)

func (w WakeController) Report(c *gin.Context) {
	var wakeUpInfo models.Wake
	err := c.ShouldBindJSON(&wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	err = wakeModel.Report(wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, wakeUpInfo)
	return
}

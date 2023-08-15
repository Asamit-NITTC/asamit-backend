package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type TargetTimeController struct{}

var targetTimeModel = new(models.TargetTimeModel)

func (t TargetTimeController) Set(c *gin.Context) {
	var registerInfo models.TargetTime
	err := c.ShouldBindJSON(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	err = targetTimeModel.Set(registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, registerInfo)
	return
}

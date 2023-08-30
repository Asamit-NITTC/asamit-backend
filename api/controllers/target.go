package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type TargetTimeController struct {
	targetTimeModel models.TargetTimeModel
}

func InitalizeTargetTimeController(t models.TargetTimeModel) *TargetTimeController {
	return &TargetTimeController{targetTimeModel: t}
}

func (t TargetTimeController) Set(c *gin.Context) {
	var registerInfo models.TargetTime
	err := c.ShouldBindJSON(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	err = t.targetTimeModel.Set(registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, registerInfo)
	return
}

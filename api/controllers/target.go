package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type TargetTimeController struct {
	targetTimeModel models.TargetTimeModel
	userModel       models.UserModel
}

func InitalizeTargetTimeController(t models.TargetTimeModel, u models.UserModel) *TargetTimeController {
	return &TargetTimeController{targetTimeModel: t, userModel: u}
}

func (t TargetTimeController) Set(c *gin.Context) {
	var requestInfo models.TargetTime
	err := c.ShouldBindJSON(&requestInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "Can't convert to json."})
		return
	}

	err = t.targetTimeModel.Set(requestInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "Failed to write to targetTimeModel."})
		return
	}
	c.JSON(http.StatusOK, requestInfo)
	return
}

func (t TargetTimeController) Get(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incomplete query parameters."})
		return
	}

	targetTime, err := t.targetTimeModel.Get(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB get error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"targetTime": targetTime})
}

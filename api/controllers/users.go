package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type UsersController struct{}

var usersModel = new(models.UsersModel)

func (u UsersController) Show(c *gin.Context) {
	uid := c.Param("uid")
	userInfo, err := usersModel.GetUserInfo(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, userInfo)
	return
}

func (u UsersController) Register(c *gin.Context) {
	var registerInfo models.Users
	err := c.ShouldBindJSON(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	err = usersModel.SetUserInfo(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, registerInfo)
	return
}

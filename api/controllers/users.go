package controllers

import (
	"net/http"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userModel models.UserModel
}

func InitializeUserController(u models.UserModel) *UserController {
	return &UserController{userModel: u}
}

func (u UserController) GetUserInfo(c *gin.Context) {
	uid := c.Param("uid")
	userInfo, err := u.userModel.GetUserInfo(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	userInfo.Sub = ""

	c.JSON(http.StatusOK, userInfo)
	return
}

func (u UserController) SignUp(c *gin.Context) {
	var registerInfo models.User
	err := c.ShouldBindJSON(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	err = u.userModel.SetUserInfo(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, registerInfo)
	return
}

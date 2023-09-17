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
	err = u.userModel.SignUpUserInfo(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, registerInfo)
	return
}

func (u UserController) ChangeUserInfo(c *gin.Context) {
	//リクエストボディからユーザー情報の取得
	var receivedUserInfo models.User
	err := c.ShouldBindJSON(&receivedUserInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	//検証関数に入れるためのUID
	uid := receivedUserInfo.UID

	subFromContext, exist := c.Get("sub")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証エラー"})
		return
	}

	//DBに保存されているSubと、LINEで認証しContextに入れたSubと一致しているか確かめる
	subFromDB, err := u.userModel.CheckExistsUser(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	if subFromContext != subFromDB {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "未認証ユーザー"})
		return
	}

	err = u.userModel.ChangeUserInfo(receivedUserInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, receivedUserInfo)
	return
}

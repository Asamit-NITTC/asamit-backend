package controllers

import (
	"net/http"

	"fmt"

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
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, err.Error(), "UID acquisition failure."})
		return
	}

	userInfo.Sub = ""

	c.JSON(http.StatusOK, userInfo)
	return
}

func (u UserController) InquirySub(c *gin.Context) {
	sub, exist := c.Get("sub")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証エラー"})
		return
	}

	convertedStringSubFromCtx := sub.(string)

	uid, err := u.userModel.GetUIDWithSub(convertedStringSubFromCtx)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "User information acquisition error."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uid": uid})
}

func (u UserController) SignUp(c *gin.Context) {
	var registerInfo models.User
	err := c.ShouldBindJSON(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "Can't convert to json."})
		return
	}

	subFromContext, exist := c.Get("sub")
	if !exist {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, "Can't get sub from context.", "Can't get sub from context."})
		return
	}

	convertedStringSubFromCtx := fmt.Sprintf("%s", subFromContext)

	existSub, err := u.userModel.CheckExistsUserWithSub(convertedStringSubFromCtx)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, err.Error(), "Can't get sub from DB."})
		return
	}

	if existSub {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, "Can't get sub from DB.", "Can't get sub from DB."})
		return
	}

	//検証が終わり次第構造体に書き込んでModelで利用できるようにする
	registerInfo.Sub = convertedStringSubFromCtx

	err = u.userModel.SignUpUserInfo(&registerInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}

	// レスポンスを返すときに見えないようにする
	registerInfo.Sub = ""
	c.JSON(http.StatusOK, registerInfo)
	return
}

func (u UserController) ChangeUserInfo(c *gin.Context) {
	//リクエストボディからユーザー情報の取得
	var receivedUserInfo models.User
	err := c.ShouldBindJSON(&receivedUserInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, err.Error(), "Can't convert to json."})
		return
	}

	err = u.userModel.ChangeUserInfo(receivedUserInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error(), "DB write error."})
		return
	}
	c.JSON(http.StatusOK, receivedUserInfo)
	return
}

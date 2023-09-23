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
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	uid := wakeUpInfo.UserUID

	// ContextからSubを取得する
	subFromContext, exist := c.Get("sub")
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証エラー"})
		return
	}

	if subFromContext == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "未認証ユーザー"})
		return
	}

	//DBからのSub取得
	subFromDB, err := w.userModel.CheckExistsUserWithUID(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	//Context(LINE認証サーバー)から受け取ったSubと、DBから取得したSubが違ったら改ざんの恐れがあるので弾く
	if subFromContext != subFromDB {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "未認証ユーザー"})
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

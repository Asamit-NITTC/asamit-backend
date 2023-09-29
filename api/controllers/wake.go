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
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusBadRequest, "Can't convert to json."})
		return
	}

	uid := wakeUpInfo.UserUID

	// ContextからSubを取得する
	subFromContext, exist := c.Get("sub")
	if !exist {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, "Can't get sub from context."})
		return
	}

	if subFromContext == "" {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, "Can't get sub from context."})
		return
	}

	//DBからのSub取得
	subFromDB, err := w.userModel.CheckExistsUserWithUID(uid)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, "Can't get sub from DB."})
		return
	}

	//Context(LINE認証サーバー)から受け取ったSubと、DBから取得したSubが違ったら改ざんの恐れがあるので弾く
	if subFromContext != subFromDB {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusUnauthorized, "The sub obtained from context and the sub obtained from db do not match."})
		return
	}

	err = w.wakeModel.Report(wakeUpInfo)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic).SetMeta(APIError{http.StatusInternalServerError, err.Error()})
		return
	}
	c.JSON(http.StatusOK, wakeUpInfo)
	return
}

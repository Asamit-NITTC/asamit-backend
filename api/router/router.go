package router

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/Asamit-NITTC/asamit-backend-test/controllers"
	"github.com/Asamit-NITTC/asamit-backend-test/middleware"
	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/Asamit-NITTC/asamit-backend-test/webstorage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, ctx context.Context, bucket *storage.BucketHandle) *gin.Engine {

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	r.Use(cors.New(config))
	r.Use(middleware.ErrorHandler())

	users := r.Group("users")
	{
		userModel := models.InitializeUserRepo(db)
		userController := controllers.InitializeUserController(userModel)
		users.GET("/:uid", userController.GetUserInfo)
		users.GET("inquiry-sub", middleware.AuthHandler(), userController.InquirySub)
		users.PUT("/update_profile", middleware.AuthHandler(), userController.ChangeUserInfo)
		users.POST("/signup", middleware.AuthHandler(), userController.SignUp)
	}

	targetTime := r.Group("target-time")
	{
		targetTimeModel := models.InitializeTargetRepo(db)
		userModel := models.InitializeUserRepo(db)
		targetTimeController := controllers.InitalizeTargetTimeController(targetTimeModel, userModel)
		targetTime.PUT("/set", middleware.AuthHandler(), targetTimeController.Set)
	}

	wake := r.Group("wake")
	{
		wakeModel := models.InitializeWakeRepo(db)
		userModel := models.InitializeUserRepo(db)
		wakeController := controllers.InitializeWakeController(wakeModel, userModel)
		wake.POST("/report", middleware.AuthHandler(), wakeController.Report)
	}

	room := r.Group("summit")
	{
		roomModel := models.InitializeRoomRepo(db)
		userModel := models.InitializeUserRepo(db)
		roomusersLinkModel := models.InitializeRoomUsersLinkRepo(db)
		approvePendingModel := models.InitializeApprovePendingRepo(db)
		roomTalkModel := models.InitializeRoomTaliRepo(db)
		cloudStorageOriginalWebModel := webstorage.InitializeCloudStorageOriginalWebRepo(ctx, bucket)
		roomController := controllers.InitailizeRoomController(roomModel, userModel, roomusersLinkModel, approvePendingModel, roomTalkModel, cloudStorageOriginalWebModel)
		room.POST("/create", middleware.AuthHandler(), roomController.Create)
		room.GET("/room-affiliation-status", middleware.AuthHandler(), roomController.CheckAffiliateAndInventionStatus)
		room.GET("/room-detail-info", middleware.AuthHandler(), roomController.GetRoomDetailInfo)
		room.POST("/record-talk", middleware.AuthHandler(), roomController.RecordTalk)
	}
	return r
}

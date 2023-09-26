package router

import (
	"github.com/Asamit-NITTC/asamit-backend-test/controllers"
	"github.com/Asamit-NITTC/asamit-backend-test/middleware"
	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {

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
	return r
}

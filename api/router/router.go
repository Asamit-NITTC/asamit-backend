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
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))
	r.Use(middleware.ErrorHandler())

	users := r.Group("users")
	{
		userModel := models.InitalizeUserRepo(db)
		userController := controllers.InitializeUserController(userModel)
		users.GET("/:uid", userController.GetUserInfo)
		users.PUT("/update_profile", middleware.AuthHandler(), userController.ChangeUserInfo)
		users.POST("/signup", userController.SignUp)
	}

	targetTime := r.Group("target-time")
	{
		targetTimeModel := models.InitializeTargetRepo(db)
		targetTimeController := controllers.InitalizeTargetTimeController(targetTimeModel)
		targetTime.PUT("/set", targetTimeController.Set)
	}

	wake := r.Group("wake")
	{
		wakeModel := models.InitalizeWakeRepo(db)
		wakeController := controllers.InitializeWakeController(wakeModel)
		wake.POST("/report", wakeController.Report)
	}
	return r
}

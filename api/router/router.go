package router

import (
	"github.com/Asamit-NITTC/asamit-backend-test/controllers"
	"github.com/Asamit-NITTC/asamit-backend-test/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))
	r.Use(middleware.ErrorHandler())

	users := r.Group("users")
	{
		u := new(controllers.UsersController)
		users.GET("/:uid", u.Show)
		users.POST("/register", u.Register)
	}

	targetTime := r.Group("target-time")
	{
		t := new(controllers.TargetTimeController)
		targetTime.PUT("/set", t.Set)
	}
	return r

}

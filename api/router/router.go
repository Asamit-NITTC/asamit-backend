package router

import (
	"github.com/Asamit-NITTC/asamit-backend-test/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	users := router.Group("users")
	{
		u := new(controllers.Users)
		users.GET("/register", u.Register)
	}

	return router

}

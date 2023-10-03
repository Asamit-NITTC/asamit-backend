package middleware

import (
	"log"
	"os"

	"github.com/Asamit-NITTC/asamit-backend-test/controllers"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {

			//デバッグ時には受け取ったエラーを生で受け取れるように
			if os.Getenv("MODE") == "DEBUG" {
				log.Println(err)
				apierror := err.Meta.(controllers.APIError)
				c.AbortWithStatusJSON(apierror.StatusCode, gin.H{
					"error": apierror.ProductionErrorMessage,
				})
			} else {
				//本番環境ではセキュリティの関係上簡単なエラーだけを表示できるように
				apierror := err.Meta.(controllers.APIError)
				c.AbortWithStatusJSON(apierror.StatusCode, gin.H{
					"error": apierror.ErrorMessage,
				})

			}
		}
	}
}

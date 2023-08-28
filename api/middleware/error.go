package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {
			log.Println(err)
			apierror := err.Meta.(APIError)
			c.AbortWithStatusJSON(apierror.StatusCode, gin.H{
				"error": apierror.ErrorMessage,
			})
		}
	}
}

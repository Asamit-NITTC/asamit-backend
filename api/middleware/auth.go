package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Asamit-NITTC/asamit-backend-test/models"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	checkI models.AuthModel
}

func InitializeAuthController(am models.AuthModel) *AuthMiddleware {
	return &AuthMiddleware{checkI: am}
}

type responseBody struct {
	Sub string `json:"sub"`
}

func (a AuthMiddleware) AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenWithBearer := c.GetHeader("Authorization")
		if tokenWithBearer == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}

		token := strings.ReplaceAll(tokenWithBearer, "Bearer ", "")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}

		v := url.Values{}
		v.Set("id_token", token)
		v.Set("client_id", os.Getenv("CLIENT_ID"))

		response, err := http.PostForm("https://api.line.me/oauth2/v2.1/verify", v)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusVariantAlsoNegotiates, gin.H{
				"code":    http.StatusVariantAlsoNegotiates,
				"message": "Connection error with LINE server",
			})
			return
		}

		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusVariantAlsoNegotiates, gin.H{
				"code":    http.StatusVariantAlsoNegotiates,
				"message": "Middleware error",
			})
			return
		}

		var responseJSON responseBody
		err = json.Unmarshal(responseBodyBytes, &responseJSON)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusVariantAlsoNegotiates, gin.H{
				"code":    http.StatusVariantAlsoNegotiates,
				"message": "JSON parse error",
			})
			return
		}

		sub := responseJSON.Sub
		c.Set("sub", sub)

		subIsValid, err := a.checkI.CheckSubIsValid(sub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusVariantAlsoNegotiates, gin.H{
				"code":    http.StatusVariantAlsoNegotiates,
				"message": "sub check error",
			})
			return
		}

		c.Set("subIsValid", subIsValid)
		c.Next()
	}
}

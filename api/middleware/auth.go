package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type responseBody struct {
	Sub string `json:"sub"`
	Err string `json:"error"`
}

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenWithBearer := c.GetHeader("Authorization")
		if tokenWithBearer == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}

		// Bearerを削除
		token := strings.ReplaceAll(tokenWithBearer, "Bearer ", "")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}

		// LINEの認証サーバーに問い合わせ
		v := url.Values{}
		v.Set("id_token", token)
		v.Set("client_id", os.Getenv("CLIENT_ID"))

		response, err := http.PostForm("https://api.line.me/oauth2/v2.1/verify", v)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Connection error with LINE server",
			})
			return
		}

		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Middleware error",
			})
			return
		}

		var responseJSON responseBody
		err = json.Unmarshal(responseBodyBytes, &responseJSON)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "JSON parse error",
			})
			return
		}

		if !(responseJSON.Err == "") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": responseJSON.Err,
			})
		}
		sub := responseJSON.Sub
		// LINE側のsub(primarykey)みたいなものをcontextに書き込んでおく
		c.Set("sub", sub)
		c.Next()
	}
}

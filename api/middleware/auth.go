package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type responseBody struct {
	Sub string `json:"sub"`
	Err string `json:"error_description"`
}

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenWithBearer := c.GetHeader("Authorization")
		if tokenWithBearer == "" {
			log.Println("Authorization header is empty.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}

		// Bearerを削除
		token := strings.ReplaceAll(tokenWithBearer, "Bearer ", "")

		if token == "" {
			log.Println("ID token is empty.")
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
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Connection error with LINE server",
			})
			return
		}

		responseBodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Middleware error",
			})
			return
		}

		var responseJSON responseBody
		err = json.Unmarshal(responseBodyBytes, &responseJSON)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "JSON parse error",
			})
			return
		}

		if !(responseJSON.Err == "") {
			log.Println("There is an error in the response json.")
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

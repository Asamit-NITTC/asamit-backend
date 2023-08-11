package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Users struct{}

func (u Users) Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Hello,World!"})
}

package handlers

import (
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InfoHandler(c *gin.Context) {
	username, ex := c.Get("user")

	if !ex {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не найден"})
		return
	}

	info, err := services.GetUserInfo(username.(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить информацию"})
		return
	}

	c.JSON(http.StatusOK, info)
}

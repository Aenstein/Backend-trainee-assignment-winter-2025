package handlers

import (
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuyHandler(c *gin.Context) {
	item := c.Param("item")
	username, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не найден"})
		return
	}

	err := services.BuyItem(username.(string), item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "покупка успешна"})
}

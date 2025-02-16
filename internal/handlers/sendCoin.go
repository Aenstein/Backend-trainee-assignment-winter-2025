package handlers

import (
	"net/http"
	"avito-shop/internal/services"
	"avito-shop/models"

	"github.com/gin-gonic/gin"
)

func SendCoinHandler(c *gin.Context) {
	var req models.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный запрос"})
		return
	}

	fromUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не найден"})
		return
	}

	err := services.SendCoins(fromUser.(string), req.ToUser, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "монеты успешно переведены"})
}

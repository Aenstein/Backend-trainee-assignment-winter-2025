package handlers

import (
	"avito-shop/internal/auth"
	"avito-shop/internal/services"
	"avito-shop/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthHandler(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный запрос"})
		return
	}

	user, err := services.GetOrCreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка работы с пользователем"})
		return
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token})
}

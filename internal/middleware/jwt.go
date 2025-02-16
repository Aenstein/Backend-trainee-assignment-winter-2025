package middleware

import (
	"avito-shop/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)


func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Токен остутствует"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			return
		}

		tokenStr := parts[1]
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не валидный токен"})
			return
		}

		c.Set("user", claims.Username)
		c.Next()
	}
}

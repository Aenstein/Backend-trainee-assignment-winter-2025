package main

import (
	"avito-shop/internal/handlers"
	"avito-shop/internal/middleware"
	"avito-shop/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	if err := services.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	router := gin.Default()

	router.POST("/api/auth", handlers.AuthHandler)

	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware())
	{
		api.GET("/info", handlers.InfoHandler)
		api.POST("/sendCoin", handlers.SendCoinHandler)
		api.GET("/buy/:item", handlers.BuyHandler)
	}

	port := "8080"
	log.Printf("Сервис запущен на порту %s", port)
	router.Run(":" + port)
}

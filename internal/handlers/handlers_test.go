package handlers_test

import (
	"avito-shop/internal/handlers"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandlerAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("empty fields", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/auth", handlers.AuthHandler)

		body := `{"username":"","password":""}`
		req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code, "Ожидается статус 400")
	})
}

func TestSendCoinHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Invalide JSON send coin", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/sendCoin", handlers.SendCoinHandler)

		req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBufferString("{invalid_json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Ожидается статус 400")
	})

	t.Run("Missing required fields", func(t *testing.T) {
		router := gin.New()
		router.POST("/api/sendCoin", func(c *gin.Context) {
			c.Set("user", "testuser")
			handlers.SendCoinHandler(c)
		})

		reqBody := `{"toUser": "anotheruser"}`
		req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Ожидается статус 400")
	})
}

func TestBuyHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Invalide user", func(t *testing.T) {
		router := gin.New()
		router.GET("/api/buy/:item", handlers.BuyHandler)

		req, _ := http.NewRequest("GET", "/api/buy/t-shirt", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Ожидается статус 401")
	})

	t.Run("Invalide merch name", func(t *testing.T) {
		router := gin.New()
		router.GET("/api/buy/:item", func(c *gin.Context) {
			c.Set("user", "testuser")
			handlers.BuyHandler(c)
		})

		req, _ := http.NewRequest("GET", "/api/buy/nonexistent_item", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Ожидается статус 400")
	})
}

func TestInfoHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Empty user", func(t *testing.T) {
		router := gin.New()
		router.GET("/api/info", handlers.InfoHandler)

		req, _ := http.NewRequest("GET", "/api/info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Ожидается статус 400")
	})
}

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"avito-shop/internal/handlers"
	"avito-shop/internal/middleware"
	"avito-shop/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Создание таблицы в контейнере для интеграционных тестов
func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "merchshop",
			"POSTGRES_USER":     "merchuser",
			"POSTGRES_PASSWORD": "merchpass",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("Ошибка запуска контейнера: %s\n", err)
		os.Exit(1)
	}

	host, err := container.Host(ctx)
	if err != nil {
		fmt.Printf("Ошибка получения host: %s\n", err)
		os.Exit(1)
	}
	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Printf("Ошибка получения порта: %s\n", err)
		os.Exit(1)
	}
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", mappedPort.Port())
	os.Setenv("DB_USER", "merchuser")
	os.Setenv("DB_PASSWORD", "merchpass")
	os.Setenv("DB_NAME", "merchshop")
	os.Setenv("JWT_SECRET", "your_secret_key")

	time.Sleep(5 * time.Second)

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		fmt.Printf("Ошибка остановки контейнера: %s\n", err)
	}

	os.Exit(code)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/auth", handlers.AuthHandler)
	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware())
	{
		api.GET("/info", handlers.InfoHandler)
		api.POST("/sendCoin", handlers.SendCoinHandler)
		api.GET("/buy/:item", handlers.BuyHandler)
	}
	return router
}

// Интеграционный тест на сценарий покупки мерча
func TestMerchPurchaseIntegration(t *testing.T) {
	err := services.InitDB()
	assert.NoError(t, err)

	router := setupRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	client := &http.Client{}

	// Регистрация пользователя
	authURL := server.URL + "/api/auth"
	userPayload := map[string]string{
		"username": "integrationUserPurchase",
		"password": "testpass",
	}
	payloadBytes, _ := json.Marshal(userPayload)
	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp map[string]string
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	resp.Body.Close()
	assert.NoError(t, err)
	token, ok := authResp["token"]
	assert.True(t, ok)
	assert.NotEmpty(t, token)

	// Совершение покупки мерча
	buyURL := server.URL + "/api/buy/t-shirt"
	req, err := http.NewRequest("GET", buyURL, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	buyResp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, buyResp.StatusCode)
	var buyResult map[string]string
	err = json.NewDecoder(buyResp.Body).Decode(&buyResult)
	buyResp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, "покупка успешна", buyResult["message"])

	// Проверка баланса, должен уменьшиться, а в инвентаре появиться t-shirt.
	infoURL := server.URL + "/api/info"
	reqInfo, err := http.NewRequest("GET", infoURL, nil)
	assert.NoError(t, err)
	reqInfo.Header.Set("Authorization", "Bearer "+token)
	infoResp, err := client.Do(reqInfo)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoResp.StatusCode)
	var infoResult map[string]interface{}
	err = json.NewDecoder(infoResp.Body).Decode(&infoResult)
	infoResp.Body.Close()
	assert.NoError(t, err)

	// Изначально у пользователя 1000 монет, t-shirt стоит 80, ожидаем 920.
	coins, ok := infoResult["coins"].(float64)
	assert.True(t, ok)
	assert.Equal(t, 920.0, coins)

	// Проверяем, что в инвентаре есть t-shirt
	inv, ok := infoResult["inventory"].([]interface{})
	assert.True(t, ok)
	found := false
	for _, item := range inv {
		itm, ok := item.(map[string]interface{})
		if ok && itm["type"] == "t-shirt" {
			quantity, ok := itm["quantity"].(float64)
			if ok && quantity >= 1 {
				found = true
				break
			}
		}
	}
	assert.True(t, found, "inventory должен содержать t-shirt")
}

// Интеграционный тест на сценарий передачи монеток другим сотрудникам
func TestCoinTransferIntegration(t *testing.T) {
	err := services.InitDB()
	assert.NoError(t, err)

	router := setupRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	client := &http.Client{}

	authURL := server.URL + "/api/auth"

	// Регистрируем пользователя A
	userAPayload := map[string]string{
		"username": "integrationUserA",
		"password": "passA",
	}
	payloadBytes, _ := json.Marshal(userAPayload)
	respA, err := http.Post(authURL, "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respA.StatusCode)
	var authRespA map[string]string
	err = json.NewDecoder(respA.Body).Decode(&authRespA)
	respA.Body.Close()
	assert.NoError(t, err)
	tokenA := authRespA["token"]
	assert.NotEmpty(t, tokenA)

	// Регистрируем пользователя B
	userBPayload := map[string]string{
		"username": "integrationUserB",
		"password": "passB",
	}
	payloadBytes, _ = json.Marshal(userBPayload)
	respB, err := http.Post(authURL, "application/json", bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respB.StatusCode)
	var authRespB map[string]string
	err = json.NewDecoder(respB.Body).Decode(&authRespB)
	respB.Body.Close()
	assert.NoError(t, err)
	tokenB := authRespB["token"]
	assert.NotEmpty(t, tokenB)

	// Пользователь A переводит 50 монет пользователю B через /api/sendCoin
	sendCoinURL := server.URL + "/api/sendCoin"
	transferPayload := map[string]interface{}{
		"toUser": "integrationUserB",
		"amount": 50,
	}
	payloadBytes, _ = json.Marshal(transferPayload)
	reqTransfer, err := http.NewRequest("POST", sendCoinURL, bytes.NewBuffer(payloadBytes))
	assert.NoError(t, err)
	reqTransfer.Header.Set("Content-Type", "application/json")
	reqTransfer.Header.Set("Authorization", "Bearer "+tokenA)
	transferResp, err := client.Do(reqTransfer)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, transferResp.StatusCode)
	var transferResult map[string]string
	err = json.NewDecoder(transferResp.Body).Decode(&transferResult)
	transferResp.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, "монеты успешно переведены", transferResult["message"])

	// Проверяем баланс пользователя A через /api/info (ожидается 950)
	infoURL := server.URL + "/api/info"
	reqInfoA, _ := http.NewRequest("GET", infoURL, nil)
	reqInfoA.Header.Set("Authorization", "Bearer "+tokenA)
	infoRespA, err := client.Do(reqInfoA)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoRespA.StatusCode)
	var infoA map[string]interface{}
	err = json.NewDecoder(infoRespA.Body).Decode(&infoA)
	infoRespA.Body.Close()
	assert.NoError(t, err)
	coinsA, ok := infoA["coins"].(float64)
	assert.True(t, ok)
	assert.Equal(t, 950.0, coinsA)

	// Проверяем баланс пользователя B через /api/info (ожидается 1050)
	reqInfoB, _ := http.NewRequest("GET", infoURL, nil)
	reqInfoB.Header.Set("Authorization", "Bearer "+tokenB)
	infoRespB, err := client.Do(reqInfoB)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoRespB.StatusCode)
	var infoB map[string]interface{}
	err = json.NewDecoder(infoRespB.Body).Decode(&infoB)
	infoRespB.Body.Close()
	assert.NoError(t, err)
	coinsB, ok := infoB["coins"].(float64)
	assert.True(t, ok)
	assert.Equal(t, 1050.0, coinsB)
}

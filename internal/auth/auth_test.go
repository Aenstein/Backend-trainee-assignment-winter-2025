package auth_test

import (
	"avito-shop/internal/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	t.Run("Generate token", func(t *testing.T) {
		token, err := auth.GenerateToken("testuser")
		assert.NoError(t, err, "должно получиться сгенерировать токен без ошибок")
		assert.NotEmpty(t, token, "токен не должен быть пустым")
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("Invalide token", func(t *testing.T) {
		claims, err := auth.ValidateToken("invalid.token.value")
		assert.Error(t, err, "должна быть ошибка при неправильном токене")
		assert.Nil(t, claims)
	})

	t.Run("Valide token", func(t *testing.T) {
		token, _ := auth.GenerateToken("testuser")

		claims, err := auth.ValidateToken(token)
		assert.NoError(t, err, "валидный токен не должен вызывать ошибку")
		assert.Equal(t, "testuser", claims.Username)
	})
}

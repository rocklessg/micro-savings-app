package tests

import (
	"testing"
	"time"

	"micro-savings-app/services"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	userID := "123456"

	token, err := services.GenerateJWT(userID)
	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	var userID = "123456"
	token, _ := services.GenerateJWT(userID)

	claims, err := services.ValidateJWT(token)
	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims["user_id"])
	assert.Greater(t, int64(claims["exp"].(float64)), time.Now().Unix())
}
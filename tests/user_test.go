package tests

import (
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"encoding/json"

	"micro-savings-app/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a request body
	requestBody := `{"name": "Abdulhafiz", "email": "test@email.com", "password": "P@ssword123"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/user/register", bytes.NewBufferString(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler
	handlers.RegisterUser(c)

	// Check the response
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a request body
	requestBody := `{"email": "test@example.com", "password": "password123"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/user/login", bytes.NewBufferString(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handlers.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "success", response["status"])
	assert.NotEmpty(t, response["token"])
}
package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"nyasah-backend/api/handlers"
	"nyasah-backend/config"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{JWTSecret: "test-secret"}

	t.Run("Register User", func(t *testing.T) {
		input := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
			"name":     "Test User",
		}

		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))

		handler := handlers.NewAuthHandler(nil, cfg)
		handler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Login User", func(t *testing.T) {
		input := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		}

		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))

		handler := handlers.NewAuthHandler(nil, cfg)
		handler.Login(c)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, response, "token")
	})
}

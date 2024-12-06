package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"nyasah-backend/api/handlers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReviewHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Create Review", func(t *testing.T) {
		input := map[string]interface{}{
			"entity_id": uuid.New(),
			"rating":    5,
			"content":   "Great product!",
		}

		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/reviews", bytes.NewBuffer(body))

		// Mock user authentication
		c.Set("user_id", uuid.New())

		handler := handlers.NewReviewHandler(nil)
		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Get Review", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reviewID := uuid.New()
		c.Request, _ = http.NewRequest("GET", "/api/reviews/"+reviewID.String(), nil)

		handler := handlers.NewReviewHandler(nil)
		handler.Get(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

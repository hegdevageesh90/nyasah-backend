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

func TestSocialProofHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Create Social Proof", func(t *testing.T) {
		input := map[string]interface{}{
			"type":      "purchase",
			"entity_id": uuid.New(),
			"content":   "User purchased this item",
		}

		body, _ := json.Marshal(input)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/social-proof", bytes.NewBuffer(body))

		// Mock user authentication
		c.Set("user_id", uuid.New())

		handler := handlers.NewSocialProofHandler(nil)
		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Get Analytics", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/social-proof/analytics", nil)

		handler := handlers.NewSocialProofHandler(nil)
		handler.GetAnalytics(c)

		var response map[string]int64
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, response, "total_proofs")
		assert.Contains(t, response, "purchase_proofs")
		assert.Contains(t, response, "review_proofs")
		assert.Contains(t, response, "view_proofs")
	})
}

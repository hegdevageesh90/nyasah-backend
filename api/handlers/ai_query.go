package handlers

import (
	"net/http"
	"nyasah-backend/models"
	"nyasah-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AIQueryHandler struct {
	db        *gorm.DB
	aiService *services.Service
}

func NewAIQueryHandler(db *gorm.DB, aiService *services.Service) *AIQueryHandler {
	return &AIQueryHandler{db: db, aiService: aiService}
}

func (h *AIQueryHandler) Query(c *gin.Context) {
	var input struct {
		Query string `json:"query" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, _ := c.Get("tenant_id")

	// Process query through AI service
	response, err := h.aiService.ProcessQuery(input.Query, tenantID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process query"})
		return
	}

	// Store query and response
	aiQuery := models.AIQuery{
		TenantID: tenantID.(uuid.UUID),
		Query:    input.Query,
		Response: response,
	}

	if err := h.db.Create(&aiQuery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save query"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":    input.Query,
		"response": response,
	})
}

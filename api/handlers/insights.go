package handlers

import (
	"net/http"
	"nyasah-backend/models"
	"nyasah-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InsightsHandler struct {
	db        *gorm.DB
	aiService *services.Service
}

func NewInsightsHandler(db *gorm.DB, aiService *services.Service) *InsightsHandler {
	return &InsightsHandler{db: db, aiService: aiService}
}

func (h *InsightsHandler) GetProductInsights(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var insights models.ProductInsights
	if err := h.db.Where("product_id = ?", productID).First(&insights).Error; err != nil {
		// Generate new insights if none exist
		insights, err = h.aiService.GenerateProductInsights(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate insights"})
			return
		}
	}

	c.JSON(http.StatusOK, insights)
}

func (h *InsightsHandler) GetRecommendations(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	recommendations, err := h.aiService.GenerateRecommendations(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

func (h *InsightsHandler) GetTrendAnalysis(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	analysis, err := h.aiService.AnalyzeTrends(tenantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze trends"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"nyasah-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	db *gorm.DB
}

func NewReviewHandler(db *gorm.DB) *ReviewHandler {
	return &ReviewHandler{db: db}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var input struct {
		ProductID string          `json:"product_id" binding:"required"`
		Rating    int             `json:"rating" binding:"required,min=1,max=5"`
		Content   string          `json:"content" binding:"required"`
		UserID    string          `json:"user_id"`
		TenantID  string          `json:"tenant_id" binding:"required"`
		MetaData  json.RawMessage `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entityID, err := uuid.Parse(input.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	review := models.Review{
		ID:        uuid.New(),
		UserID:    input.UserID,
		TenantID:  input.TenantID,
		EntityID:  entityID,
		Rating:    input.Rating,
		Content:   input.Content,
		Verified:  true,
		Metadata:  input.MetaData,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *ReviewHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID is required"})
		return
	}

	var reviews []models.Review
	if err := h.db.Where("tenant_id = ?", tenantID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var review models.Review
	if err := h.db.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

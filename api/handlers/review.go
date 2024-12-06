package handlers

import (
	"net/http"
	"nyasah-backend/models"

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
		ProductID uuid.UUID `json:"product_id" binding:"required"`
		Rating    int       `json:"rating" binding:"required,min=1,max=5"`
		Content   string    `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	review := models.Review{
		UserID:   userID.(uuid.UUID),
		EntityID: input.ProductID,
		Rating:   input.Rating,
		Content:  input.Content,
		Verified: true,
	}

	if err := h.db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *ReviewHandler) List(c *gin.Context) {
	var reviews []models.Review
	if err := h.db.Find(&reviews).Error; err != nil {
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

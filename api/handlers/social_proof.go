package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"nyasah/models"
)

type SocialProofHandler struct {
	db *gorm.DB
}

func NewSocialProofHandler(db *gorm.DB) *SocialProofHandler {
	return &SocialProofHandler{db: db}
}

func (h *SocialProofHandler) Create(c *gin.Context) {
	var input struct {
		Type      string    `json:"type" binding:"required"`
		ProductID uuid.UUID `json:"product_id" binding:"required"`
		Content   string    `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	proof := models.SocialProof{
		Type:      input.Type,
		ProductID: input.ProductID,
		UserID:    userID.(uuid.UUID),
		Content:   input.Content,
	}

	if err := h.db.Create(&proof).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social proof"})
		return
	}

	c.JSON(http.StatusCreated, proof)
}

func (h *SocialProofHandler) List(c *gin.Context) {
	var proofs []models.SocialProof
	if err := h.db.Find(&proofs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social proofs"})
		return
	}

	c.JSON(http.StatusOK, proofs)
}

func (h *SocialProofHandler) GetAnalytics(c *gin.Context) {
	var stats struct {
		TotalProofs     int64 `json:"total_proofs"`
		PurchaseProofs  int64 `json:"purchase_proofs"`
		ReviewProofs    int64 `json:"review_proofs"`
		ViewProofs      int64 `json:"view_proofs"`
	}

	h.db.Model(&models.SocialProof{}).Count(&stats.TotalProofs)
	h.db.Model(&models.SocialProof{}).Where("type = ?", "purchase").Count(&stats.PurchaseProofs)
	h.db.Model(&models.SocialProof{}).Where("type = ?", "review").Count(&stats.ReviewProofs)
	h.db.Model(&models.SocialProof{}).Where("type = ?", "view").Count(&stats.ViewProofs)

	c.JSON(http.StatusOK, stats)
}
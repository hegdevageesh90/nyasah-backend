package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"nyasah/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantHandler struct {
	db *gorm.DB
}

func NewTenantHandler(db *gorm.DB) *TenantHandler {
	return &TenantHandler{db: db}
}

func generateApiKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (h *TenantHandler) Create(c *gin.Context) {
	var input struct {
		Name     string                 `json:"name" binding:"required"`
		Domain   string                 `json:"domain" binding:"required"`
		Type     string                 `json:"type" binding:"required"`
		Settings map[string]interface{} `json:"settings"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey, err := generateApiKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	tenant := models.Tenant{
		Name:     input.Name,
		Domain:   input.Domain,
		Type:     input.Type,
		ApiKey:   apiKey,
		Active:   true,
		Settings: models.JSON(input.Settings),
	}

	if err := h.db.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tenant"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      tenant.ID,
		"api_key": tenant.ApiKey,
		"message": "Tenant created successfully",
	})
}

func (h *TenantHandler) Get(c *gin.Context) {
	tenantID := c.Param("id")
	var tenant models.Tenant

	if err := h.db.First(&tenant, "id = ?", tenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) Update(c *gin.Context) {
	tenantID := c.Param("id")
	var input struct {
		Name     string                 `json:"name"`
		Settings map[string]interface{} `json:"settings"`
		Active   *bool                  `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tenant models.Tenant
	if err := h.db.First(&tenant, "id = ?", tenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Settings != nil {
		updates["settings"] = input.Settings
	}
	if input.Active != nil {
		updates["active"] = *input.Active
	}

	if err := h.db.Model(&tenant).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tenant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant updated successfully"})
}

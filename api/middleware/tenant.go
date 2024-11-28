package middleware

import (
	"net/http"
	"nyasah/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TenantMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		var tenant models.Tenant
		if err := db.Where("api_key = ? AND active = ?", apiKey, true).First(&tenant).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or inactive API key"})
			c.Abort()
			return
		}

		c.Set("tenant_id", tenant.ID)
		c.Set("tenant_type", tenant.Type)
		c.Next()
	}
}

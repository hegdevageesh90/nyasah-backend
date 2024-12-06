package models_test

import (
	"nyasah-backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenantModel(t *testing.T) {
	t.Run("BeforeCreate sets UUID", func(t *testing.T) {
		tenant := &models.Tenant{
			Name:   "Test Tenant",
			Domain: "test.com",
			Type:   "ecommerce",
		}

		err := tenant.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, tenant.ID)
	})
}

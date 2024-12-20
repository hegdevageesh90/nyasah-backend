package database

import (
	"nyasah-backend/config"
	"nyasah-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Initialize(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.Entity{},
		&models.Review{},
		&models.SocialProof{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

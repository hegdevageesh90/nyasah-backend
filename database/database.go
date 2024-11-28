package database

import (
	"nyasah/config"
	"nyasah/models"
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
		&models.User{},
		&models.Review{},
		&models.Product{},
		&models.SocialProof{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
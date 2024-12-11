package utils

import (
	"nyasah-backend/models"
	"time"

	"gorm.io/gorm"
)

// db variable to hold the database connection
var db *gorm.DB

// InitializeDB initializes the database connection
func InitializeDB(database *gorm.DB) {
	db = database
}

type TimeFrame struct {
	Reviews []models.Review
	Proofs  []models.SocialProof
	Start   time.Time
	End     time.Time
}

func GetTimeFrames() []TimeFrame {
	now := time.Now()
	frames := make([]TimeFrame, 12) // Last 12 weeks

	for i := range frames {
		frames[i].End = now.AddDate(0, 0, -7*i)
		frames[i].Start = frames[i].End.AddDate(0, 0, -7)
	}

	return frames
}

func GetReviewsInTimeFrame(tenantID string, start, end time.Time) ([]models.Review, error) {
	var reviews []models.Review
	err := db.Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, start, end).Find(&reviews).Error
	return reviews, err
}

func GetProofsInTimeFrame(tenantID string, start, end time.Time) ([]models.SocialProof, error) {
	var proofs []models.SocialProof
	err := db.Where("tenant_id = ? AND created_at BETWEEN ? AND ?", tenantID, start, end).Find(&proofs).Error
	return proofs, err
}

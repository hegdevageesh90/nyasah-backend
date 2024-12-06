package services

import (
	"nyasah-backend/models"
	"nyasah-backend/services/ai/analyzers"
	"nyasah-backend/services/ai/recommenders"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	analyzer    *analyzers.ContentAnalyzer
	recommender *recommenders.Recommender
}

func NewAIService(db *gorm.DB) *Service {
	apiKey := os.Getenv("OPENAI_API_KEY")
	return &Service{
		db:          db,
		analyzer:    analyzers.NewContentAnalyzer(apiKey),
		recommender: recommenders.NewRecommender(db, apiKey),
	}
}

func (s *Service) ProcessQuery(query string, tenantID string) (string, error) {
	// Process natural language query using the content analyzer
	response, err := s.analyzer.ProcessQuery(query, tenantID)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (s *Service) GenerateProductInsights(productID uuid.UUID) (models.ProductInsights, error) {
	// Analyze reviews and social proofs for the product
	var reviews []models.Review
	var socialProofs []models.SocialProof

	s.db.Where("product_id = ?", productID).Find(&reviews)
	s.db.Where("product_id = ?", productID).Find(&socialProofs)

	insights := models.ProductInsights{
		ProductID: productID,
	}

	// Generate sentiment trends
	insights.SentimentTrend = s.analyzer.AnalyzeSentimentTrend(reviews)

	// Extract top keywords
	insights.TopKeywords = s.analyzer.ExtractKeywords(reviews)

	// Calculate engagement score
	insights.EngagementScore = s.analyzer.CalculateEngagementScore(reviews, socialProofs)

	// Generate recommendations
	insights.RecommendedActions = s.recommender.GenerateActions(productID)

	return insights, nil
}

func (s *Service) GenerateRecommendations(tenantID uuid.UUID) ([]models.AIRecommendation, error) {
	return s.recommender.GenerateRecommendations(tenantID)
}

func (s *Service) AnalyzeTrends(tenantID uuid.UUID) (map[string]interface{}, error) {
	var reviews []models.Review
	var socialProofs []models.SocialProof

	s.db.Where("tenant_id = ?", tenantID).Find(&reviews)
	s.db.Where("tenant_id = ?", tenantID).Find(&socialProofs)

	return s.analyzer.AnalyzeTrends(reviews, socialProofs)
}

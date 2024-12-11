package services

import (
	"nyasah-backend/config"
	"nyasah-backend/models"
	"nyasah-backend/services/ai/analyzers"
	"nyasah-backend/services/ai/factory"
	"nyasah-backend/services/ai/recommenders"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	analyzer    *analyzers.ContentAnalyzer
	recommender *recommenders.Recommender
	config      *config.Config
}

func NewAIService(db *gorm.DB, config *config.Config) *Service {
	providerConfig := map[string]string{
		"api_key":    os.Getenv("OPENAI_API_KEY"),
		"model":      config.Model,
		"server_url": os.Getenv("LLAMA_SERVER_URL"),
	}

	provider, err := factory.CreateProvider(config.Provider, providerConfig)
	if err != nil {
		// Fallback to OpenAI if specified provider fails
		config.Provider = factory.OpenAI
		provider, _ = factory.CreateProvider(factory.OpenAI, providerConfig)
	}

	return &Service{
		db:          db,
		analyzer:    analyzers.NewContentAnalyzer(provider),
		recommender: recommenders.NewRecommender(db, provider),
		config:      config,
	}
}

func (s *Service) UpdateConfig(config *config.Config) error {
	providerConfig := map[string]string{
		"api_key":    os.Getenv("OPENAI_API_KEY"),
		"model":      config.Model,
		"server_url": os.Getenv("LLAMA_SERVER_URL"),
	}

	provider, err := factory.CreateProvider(config.Provider, providerConfig)
	if err != nil {
		return err
	}

	s.analyzer = analyzers.NewContentAnalyzer(provider)
	s.recommender = recommenders.NewRecommender(s.db, provider)
	s.config = config

	return nil
}

func (s *Service) ProcessQuery(query string, tenantID string) (string, error) {
	return s.analyzer.ProcessQuery(query, tenantID)
}

func (s *Service) GenerateProductInsights(productID uuid.UUID) (models.ProductInsights, error) {
	return s.recommender.GenerateInsights(productID)
}

func (s *Service) GenerateRecommendations(tenantID uuid.UUID) ([]models.AIRecommendation, error) {
	return s.recommender.GenerateRecommendations(tenantID)
}

func (s *Service) AnalyzeTrends(tenantID uuid.UUID) (map[string]interface{}, error) {
	return s.analyzer.AnalyzeTrends(tenantID.String())
}

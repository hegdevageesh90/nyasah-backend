package recommenders

import (
	"context"
	"fmt"
	"nyasah-backend/models"
	"nyasah-backend/services/ai/utils"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type Recommender struct {
	db     *gorm.DB
	client *openai.Client
}

func NewRecommender(db *gorm.DB, apiKey string) *Recommender {
	return &Recommender{
		db:     db,
		client: openai.NewClient(apiKey),
	}
}

func (r *Recommender) GenerateActions(productID uuid.UUID) []string {
	var product models.Product
	var reviews []models.Review
	var proofs []models.SocialProof

	r.db.First(&product, productID)
	r.db.Where("product_id = ?", productID).Find(&reviews)
	r.db.Where("product_id = ?", productID).Find(&proofs)

	insights := utils.GenerateProductInsights(product, reviews, proofs)

	prompt := `Based on the following product insights, suggest specific actions to improve social proof and engagement:
	
	Product: ` + product.Name + `
	Average Rating: ` + fmt.Sprintf("%.2f", insights.AverageRating) + `
	Review Count: ` + fmt.Sprintf("%d", len(reviews)) + `
	Engagement Rate: ` + fmt.Sprintf("%.2f", insights.EngagementRate) + `
	
	Provide 3-5 specific, actionable recommendations:`

	resp, err := r.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return []string{"Highlight positive reviews", "Add customer photos", "Display purchase notifications"}
	}

	return utils.ParseRecommendations(resp.Choices[0].Message.Content)
}

func (r *Recommender) GenerateRecommendations(tenantID uuid.UUID) ([]models.AIRecommendation, error) {
	var recommendations []models.AIRecommendation

	// Analyze patterns across successful products
	patterns := r.analyzeSuccessPatterns(tenantID)

	// Generate specific recommendations based on patterns
	for _, pattern := range patterns {
		confidence := utils.CalculateConfidence(pattern.SampleSize, pattern.SuccessRate)

		recommendation := models.AIRecommendation{
			TenantID:   tenantID,
			Type:       pattern.Type,
			Suggestion: pattern.GenerateRecommendation(),
			Confidence: confidence,
		}

		recommendations = append(recommendations, recommendation)
	}

	return recommendations, nil
}

func (r *Recommender) analyzeSuccessPatterns(tenantID uuid.UUID) []utils.Pattern {
	var patterns []utils.Pattern

	// Analyze content patterns
	contentPattern := r.analyzeContentPattern(tenantID)
	patterns = append(patterns, contentPattern)

	// Analyze timing patterns
	timingPattern := r.analyzeTimingPattern(tenantID)
	patterns = append(patterns, timingPattern)

	// Analyze placement patterns
	placementPattern := r.analyzePlacementPattern(tenantID)
	patterns = append(patterns, placementPattern)

	return patterns
}

func (r *Recommender) analyzeContentPattern(tenantID uuid.UUID) utils.Pattern {
	var proofs []models.SocialProof
	r.db.Joins("JOIN proof_performances ON social_proofs.id = proof_performances.proof_id").
		Where("tenant_id = ? AND proof_performances.engagement_rate > ?", tenantID, 0.7).
		Find(&proofs)

	return utils.AnalyzeContentPattern(proofs)
}

func (r *Recommender) analyzeTimingPattern(tenantID uuid.UUID) utils.Pattern {
	var proofs []models.SocialProof
	r.db.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(1000).
		Find(&proofs)

	return utils.AnalyzeTimingPattern(proofs)
}

func (r *Recommender) analyzePlacementPattern(tenantID uuid.UUID) utils.Pattern {
	var proofs []models.SocialProof
	r.db.Joins("JOIN proof_performances ON social_proofs.id = proof_performances.proof_id").
		Where("tenant_id = ?", tenantID).
		Find(&proofs)

	return utils.AnalyzePlacementPattern(proofs)
}

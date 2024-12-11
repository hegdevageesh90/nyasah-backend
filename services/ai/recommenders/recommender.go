package recommenders

import (
	"fmt"
	"nyasah-backend/models"
	"nyasah-backend/services/ai/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"nyasah-backend/services/ai/providers"
)

type Recommender struct {
	db       *gorm.DB
	provider providers.Provider
}

func NewRecommender(db *gorm.DB, provider providers.Provider) *Recommender {
	return &Recommender{
		db:       db,
		provider: provider,
	}
}

func (r *Recommender) GenerateActions(productID uuid.UUID) []string {
	var product models.Product
	var reviews []models.Review
	var proofs []models.SocialProof

	// Fetching the required data from the database
	r.db.First(&product, productID)
	r.db.Where("product_id = ?", productID).Find(&reviews)
	r.db.Where("product_id = ?", productID).Find(&proofs)

	// Generate product insights using utility
	insights := utils.GenerateProductInsights(product, reviews, proofs)

	// Constructing the prompt for the AI provider
	prompt := fmt.Sprintf(`Based on the following product insights, suggest specific actions to improve social proof and engagement:
	
	Product: %s
	Average Rating: %.2f
	Review Count: %d
	Engagement Rate: %.2f
	
	Provide 3-5 specific, actionable recommendations:`,
		product.Name,
		insights.AverageRating,
		len(reviews),
		insights.EngagementRate,
	)

	// Call the ProcessQuery method of the provider
	response, err := r.provider.ProcessQuery(prompt)
	if err != nil {
		// Return default recommendations in case of error
		return []string{"Highlight positive reviews", "Add customer photos", "Display purchase notifications"}
	}

	// Parse and return recommendations using utility
	return utils.ParseRecommendations(response)
}

// GenerateInsights fetches data for the specified product and analyzes it to produce insights.
func (r *Recommender) GenerateInsights(productID uuid.UUID) (models.ProductInsights, error) {
	var product models.Product
	var reviews []models.Review
	var proofs []models.SocialProof

	// Fetch product details
	if err := r.db.First(&product, productID).Error; err != nil {
		return models.ProductInsights{}, fmt.Errorf("failed to fetch product: %w", err)
	}

	// Fetch associated reviews
	if err := r.db.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		return models.ProductInsights{}, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	// Fetch associated social proofs
	if err := r.db.Where("product_id = ?", productID).Find(&proofs).Error; err != nil {
		return models.ProductInsights{}, fmt.Errorf("failed to fetch social proofs: %w", err)
	}

	// Calculate insights using the utility function
	insights := utils.GenerateProductInsights(product, reviews, proofs)

	// Return the analyzed insights
	return insights, nil
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

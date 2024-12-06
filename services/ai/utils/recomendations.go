package utils

import (
	"nyasah-backend/models"
	"sort"
	"strings"
)

// ParseRecommendations parses the AI response into a list of recommendations
func ParseRecommendations(aiResponse string) []string {
	// Split response by newlines or numbers
	lines := strings.Split(aiResponse, "\n")
	recommendations := make([]string, 0)

	for _, line := range lines {
		// Clean the line
		line = strings.TrimSpace(line)

		// Remove numbered bullets (1., 2., etc.)
		line = strings.TrimLeft(line, "0123456789. ")

		// Skip empty lines
		if line == "" {
			continue
		}

		recommendations = append(recommendations, line)
	}

	return recommendations
}

// CalculateConfidence calculates the confidence score for a recommendation
func CalculateConfidence(sampleSize int, successRate float64) float64 {
	// Base confidence on sample size and success rate
	sampleWeight := float64(sampleSize) / 1000 // Normalize sample size
	if sampleWeight > 1 {
		sampleWeight = 1
	}

	// Combine sample size weight with success rate
	confidence := (sampleWeight * 0.4) + (successRate * 0.6)

	// Ensure confidence is between 0 and 1
	if confidence > 1 {
		confidence = 1
	}

	return confidence
}

// FindTopPerformers identifies top performing content
func FindTopPerformers(reviews []models.Review, proofs []models.SocialProof) []map[string]interface{} {
	type performer struct {
		Type       string
		ID         string
		Score      float64
		Content    string
		Engagement float64
		CreatedAt  string
	}

	var performers []performer

	// Add reviews to performers
	for _, review := range reviews {
		score := CalculateReviewEngagement(review)
		performers = append(performers, performer{
			Type:       "review",
			ID:         review.ID.String(),
			Score:      score,
			Content:    review.Content,
			Engagement: float64(review.Engagement.Likes+review.Engagement.Shares) / float64(review.Engagement.Views),
			CreatedAt:  review.CreatedAt.Format("2006-01-02"),
		})
	}

	// Add social proofs to performers
	for _, proof := range proofs {
		score := CalculateProofEngagement(proof)
		performers = append(performers, performer{
			Type:       proof.Type,
			ID:         proof.ID.String(),
			Score:      score,
			Content:    proof.Content,
			Engagement: proof.Performance.EngagementRate,
			CreatedAt:  proof.CreatedAt.Format("2006-01-02"),
		})
	}

	// Sort performers by score
	sort.Slice(performers, func(i, j int) bool {
		return performers[i].Score > performers[j].Score
	})

	// Convert top 10 performers to map
	result := make([]map[string]interface{}, 0)
	for i := 0; i < len(performers) && i < 10; i++ {
		result = append(result, map[string]interface{}{
			"type":       performers[i].Type,
			"id":         performers[i].ID,
			"score":      performers[i].Score,
			"content":    performers[i].Content,
			"engagement": performers[i].Engagement,
			"created_at": performers[i].CreatedAt,
		})
	}

	return result
}

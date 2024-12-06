package utils

import (
	"nyasah-backend/models"
	"sort"
)

type ProductInsightData struct {
	AverageRating  float64
	EngagementRate float64
	SentimentScore float64
	TopKeywords    []string
}

func GenerateProductInsights(product models.Product, reviews []models.Review, proofs []models.SocialProof) ProductInsightData {
	insights := ProductInsightData{}

	// Calculate average rating
	var totalRating float64
	for _, review := range reviews {
		totalRating += float64(review.Rating)
	}
	if len(reviews) > 0 {
		insights.AverageRating = totalRating / float64(len(reviews))
	}

	// Calculate engagement rate
	insights.EngagementRate = calculateOverallEngagement(reviews, proofs)

	// Calculate sentiment score
	insights.SentimentScore = calculateOverallSentiment(reviews)

	// Extract top keywords
	insights.TopKeywords = extractTopKeywords(reviews)

	return insights
}

func calculateOverallEngagement(reviews []models.Review, proofs []models.SocialProof) float64 {
	var totalEngagement float64
	count := 0

	for _, review := range reviews {
		if review.Engagement.Views > 0 {
			engagement := float64(review.Engagement.Likes+review.Engagement.Shares) / float64(review.Engagement.Views)
			totalEngagement += engagement
			count++
		}
	}

	for _, proof := range proofs {
		if proof.Performance.Views > 0 {
			engagement := float64(proof.Performance.Conversions) / float64(proof.Performance.Views)
			totalEngagement += engagement
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return totalEngagement / float64(count)
}

func calculateOverallSentiment(reviews []models.Review) float64 {
	if len(reviews) == 0 {
		return 0
	}

	var totalSentiment float64
	for _, review := range reviews {
		totalSentiment += review.Sentiment
	}

	return totalSentiment / float64(len(reviews))
}

func extractTopKeywords(reviews []models.Review) []string {
	keywordFreq := make(map[string]int)

	for _, review := range reviews {
		for _, keyword := range review.Keywords {
			keywordFreq[keyword]++
		}
	}

	// Convert map to slice for sorting
	type keywordCount struct {
		word  string
		count int
	}

	var keywords []keywordCount
	for word, count := range keywordFreq {
		keywords = append(keywords, keywordCount{word, count})
	}

	// Sort by frequency
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].count > keywords[j].count
	})

	// Return top 5 keywords
	result := make([]string, 0)
	for i := 0; i < len(keywords) && i < 5; i++ {
		result = append(result, keywords[i].word)
	}

	return result
}

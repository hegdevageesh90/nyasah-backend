package analyzers

import (
	"nyasah-backend/models"
	"nyasah-backend/services/ai/providers"
	"nyasah-backend/services/ai/utils"
)

type SentimentAnalyzer struct {
	provider providers.Provider
}

func NewSentimentAnalyzer(provider providers.Provider) *SentimentAnalyzer {
	return &SentimentAnalyzer{
		provider: provider,
	}
}

func (sa *SentimentAnalyzer) AnalyzeSentiment(text string) (float64, error) {
	return sa.provider.AnalyzeSentiment(text)
}

func (sa *SentimentAnalyzer) BatchAnalyzeSentiment(reviews []models.Review) ([]float64, error) {
	scores := make([]float64, len(reviews))
	for i, review := range reviews {
		score, err := sa.AnalyzeSentiment(review.Content)
		if err != nil {
			return nil, err
		}
		scores[i] = score
	}
	return scores, nil
}

func (sa *SentimentAnalyzer) AnalyzeTrends(tenantID string) (map[string]interface{}, error) {
	timeFrames := utils.GetTimeFrames()
	trends := make([]float64, len(timeFrames))

	for i, frame := range timeFrames {
		reviews, err := utils.GetReviewsInTimeFrame(tenantID, frame.Start, frame.End)
		if err != nil {
			return nil, err
		}

		scores, err := sa.BatchAnalyzeSentiment(reviews)
		if err != nil {
			return nil, err
		}

		trends[i] = utils.CalculateAverageScore(scores)
	}

	return map[string]interface{}{
		"trends": trends,
		"frames": timeFrames,
	}, nil
}

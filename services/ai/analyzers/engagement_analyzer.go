package analyzers

import (
	"nyasah-backend/models"
	"nyasah-backend/services/ai/providers"
	"nyasah-backend/services/ai/utils"
)

type EngagementAnalyzer struct {
	provider providers.Provider
}

func NewEngagementAnalyzer(provider providers.Provider) *EngagementAnalyzer {
	return &EngagementAnalyzer{
		provider: provider,
	}
}

func (ea *EngagementAnalyzer) AnalyzeEngagement(review models.Review, proof models.SocialProof) float64 {
	reviewScore := utils.CalculateReviewEngagement(review)
	proofScore := utils.CalculateProofEngagement(proof)
	return (reviewScore + proofScore) / 2
}

func (ea *EngagementAnalyzer) AnalyzeTrends(tenantID string) (map[string]interface{}, error) {
	timeFrames := utils.GetTimeFrames()
	trends := make([]float64, len(timeFrames))

	for i, frame := range timeFrames {
		reviews, err := utils.GetReviewsInTimeFrame(tenantID, frame.Start, frame.End)
		if err != nil {
			return nil, err
		}

		proofs, err := utils.GetProofsInTimeFrame(tenantID, frame.Start, frame.End)
		if err != nil {
			return nil, err
		}

		trends[i] = utils.CalculateAverageEngagement(reviews, proofs)
	}

	return map[string]interface{}{
		"trends": trends,
		"frames": timeFrames,
	}, nil
}

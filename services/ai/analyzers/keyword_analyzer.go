package analyzers

import (
	"nyasah-backend/services/ai/providers"
	"nyasah-backend/services/ai/utils"
)

type KeywordAnalyzer struct {
	provider providers.Provider
}

func NewKeywordAnalyzer(provider providers.Provider) *KeywordAnalyzer {
	return &KeywordAnalyzer{
		provider: provider,
	}
}

func (ka *KeywordAnalyzer) ExtractKeywords(text string) ([]string, error) {
	prompt := `Extract key phrases and topics from the following text:
	
	Text: """` + text + `"""
	
	Return only the key phrases, separated by commas:`

	response, err := ka.provider.ProcessQuery(prompt)
	if err != nil {
		return nil, err
	}

	return utils.ParseKeywords(response), nil
}

func (ka *KeywordAnalyzer) AnalyzeTrends(tenantID string) (map[string]interface{}, error) {
	timeFrames := utils.GetTimeFrames()
	trends := make([]map[string]int, len(timeFrames))

	for i, frame := range timeFrames {
		reviews, err := utils.GetReviewsInTimeFrame(tenantID, frame.Start, frame.End)
		if err != nil {
			return nil, err
		}

		keywordFreq := make(map[string]int)
		for _, review := range reviews {
			keywords, err := ka.ExtractKeywords(review.Content)
			if err != nil {
				continue
			}

			for _, keyword := range keywords {
				keywordFreq[keyword]++
			}
		}

		trends[i] = keywordFreq
	}

	return map[string]interface{}{
		"trends": trends,
		"frames": timeFrames,
	}, nil
}

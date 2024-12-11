package analyzers

import (
	"nyasah-backend/services/ai/providers"
)

type ContentAnalyzer struct {
	provider providers.Provider
}

func NewContentAnalyzer(provider providers.Provider) *ContentAnalyzer {
	return &ContentAnalyzer{
		provider: provider,
	}
}

func (a *ContentAnalyzer) ProcessQuery(query, tenantID string) (string, error) {
	prompt := `As an AI assistant for an e-commerce social proof platform, answer the following query:
	
	Query: """` + query + `"""
	
	Provide a clear, concise, and helpful response based on the available data and best practices.`

	return a.provider.ProcessQuery(prompt)
}

func (a *ContentAnalyzer) AnalyzeTrends(tenantID string) (map[string]interface{}, error) {
	// Delegate to specialized analyzers
	sentimentAnalyzer := NewSentimentAnalyzer(a.provider)
	engagementAnalyzer := NewEngagementAnalyzer(a.provider)
	keywordAnalyzer := NewKeywordAnalyzer(a.provider)

	sentimentTrends, err := sentimentAnalyzer.AnalyzeTrends(tenantID)
	if err != nil {
		return nil, err
	}

	engagementTrends, err := engagementAnalyzer.AnalyzeTrends(tenantID)
	if err != nil {
		return nil, err
	}

	keywordTrends, err := keywordAnalyzer.AnalyzeTrends(tenantID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"sentiment_trends":  sentimentTrends,
		"engagement_trends": engagementTrends,
		"keyword_trends":    keywordTrends,
	}, nil
}

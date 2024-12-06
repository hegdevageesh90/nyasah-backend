package analyzers

import (
	"context"
	"nyasah-backend/models"
	"nyasah-backend/services/ai/utils"
	"sort"
	"time"

	"github.com/sashabaranov/go-openai"
)

type ContentAnalyzer struct {
	client            *openai.Client
	sentimentAnalyzer *SentimentAnalyzer
}

func NewContentAnalyzer(apiKey string) *ContentAnalyzer {
	return &ContentAnalyzer{
		client:            openai.NewClient(apiKey),
		sentimentAnalyzer: NewSentimentAnalyzer(apiKey),
	}
}

func (a *ContentAnalyzer) ProcessQuery(query, tenantID string) (string, error) {
	prompt := `As an AI assistant for an e-commerce social proof platform, answer the following query:
	
	Query: """` + query + `"""
	
	Provide a clear, concise, and helpful response based on the available data and best practices.`

	resp, err := a.client.CreateChatCompletion(
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
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (a *ContentAnalyzer) AnalyzeSentimentTrend(reviews []models.Review) []float64 {
	scores, err := a.sentimentAnalyzer.BatchAnalyzeSentiment(reviews)
	if err != nil {
		return make([]float64, 0)
	}
	return scores
}

func (a *ContentAnalyzer) ExtractKeywords(reviews []models.Review) []string {
	keywordFreq := make(map[string]int)

	for _, review := range reviews {
		phrases, err := a.sentimentAnalyzer.ExtractKeyPhrases(review.Content)
		if err != nil {
			continue
		}

		for _, phrase := range phrases {
			keywordFreq[phrase]++
		}
	}

	// Sort keywords by frequency
	type kw struct {
		word  string
		count int
	}

	var keywords []kw
	for word, count := range keywordFreq {
		keywords = append(keywords, kw{word, count})
	}

	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].count > keywords[j].count
	})

	// Return top 10 keywords
	result := make([]string, 0)
	for i := 0; i < len(keywords) && i < 10; i++ {
		result = append(result, keywords[i].word)
	}

	return result
}

func (a *ContentAnalyzer) CalculateEngagementScore(reviews []models.Review, proofs []models.SocialProof) float64 {
	if len(reviews) == 0 && len(proofs) == 0 {
		return 0
	}

	var totalScore float64

	// Review engagement
	for _, review := range reviews {
		score := utils.CalculateReviewEngagement(review)
		totalScore += score
	}

	// Social proof engagement
	for _, proof := range proofs {
		score := utils.CalculateProofEngagement(proof)
		totalScore += score
	}

	totalItems := float64(len(reviews) + len(proofs))
	return totalScore / totalItems
}

func (a *ContentAnalyzer) AnalyzeTrends(reviews []models.Review, proofs []models.SocialProof) (map[string]interface{}, error) {
	// Group reviews and proofs by time periods
	timeFrames := utils.GroupByTimeFrames(reviews, proofs)

	// Calculate trends
	sentimentTrend := make([]float64, len(timeFrames))
	engagementTrend := make([]float64, len(timeFrames))

	for i, frame := range timeFrames {
		sentimentTrend[i] = utils.CalculateAverageSentiment(frame.Reviews)
		engagementTrend[i] = utils.CalculateAverageEngagement(frame.Reviews, frame.Proofs)
	}

	// Find top performers
	topPerformers := utils.FindTopPerformers(reviews, proofs)

	return map[string]interface{}{
		"sentiment_trend":  sentimentTrend,
		"engagement_trend": engagementTrend,
		"top_performers":   topPerformers,
		"analysis_time":    time.Now(),
	}, nil
}

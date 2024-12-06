package analyzers

import (
	"context"
	"fmt"
	"nyasah-backend/models"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type SentimentAnalyzer struct {
	client *openai.Client
}

func NewSentimentAnalyzer(apiKey string) *SentimentAnalyzer {
	return &SentimentAnalyzer{
		client: openai.NewClient(apiKey),
	}
}

func (sa *SentimentAnalyzer) AnalyzeSentiment(text string) (float64, error) {
	prompt := `Analyze the sentiment of the following text and return a score between -1 (very negative) and 1 (very positive):
	
	Text: """` + text + `"""
	
	Score:`

	resp, err := sa.client.CreateChatCompletion(
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
		return 0, err
	}

	score := 0.0
	_, err = fmt.Sscanf(resp.Choices[0].Message.Content, "%f", &score)
	if err != nil {
		return 0, err
	}

	return score, nil
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

func (sa *SentimentAnalyzer) ExtractKeyPhrases(text string) ([]string, error) {
	prompt := `Extract the key phrases from the following text that represent the main topics or sentiments:
	
	Text: """` + text + `"""
	
	Key phrases (comma-separated):`

	resp, err := sa.client.CreateChatCompletion(
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
		return nil, err
	}

	phrases := strings.Split(resp.Choices[0].Message.Content, ",")
	for i := range phrases {
		phrases[i] = strings.TrimSpace(phrases[i])
	}

	return phrases, nil
}

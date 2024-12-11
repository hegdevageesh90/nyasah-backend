package providers

type Provider interface {
	ProcessQuery(query string) (string, error)
	AnalyzeSentiment(text string) (float64, error)
	GenerateText(prompt string, maxTokens int, temperature float64) (string, error)
}

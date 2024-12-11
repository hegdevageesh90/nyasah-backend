package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LlamaProvider struct {
	serverURL string
}

func NewLlamaProvider(serverURL string) *LlamaProvider {
	return &LlamaProvider{
		serverURL: serverURL,
	}
}

func (p *LlamaProvider) ProcessQuery(query string) (string, error) {
	payload := map[string]interface{}{
		"prompt":      query,
		"max_tokens":  1000,
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(p.serverURL+"/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Text, nil
}

func (p *LlamaProvider) AnalyzeSentiment(text string) (float64, error) {
	prompt := fmt.Sprintf(`Analyze the sentiment of the following text and return a score between -1 (very negative) and 1 (very positive):
	
	Text: """%s"""
	
	Score:`, text)

	payload := map[string]interface{}{
		"prompt":      prompt,
		"max_tokens":  100,
		"temperature": 0.3,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(p.serverURL+"/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	var score float64
	_, err = fmt.Sscanf(result.Text, "%f", &score)
	return score, err
}

func (p *LlamaProvider) GenerateText(prompt string, maxTokens int, temperature float64) (string, error) {
	payload := map[string]interface{}{
		"prompt":      prompt,
		"max_tokens":  maxTokens,
		"temperature": temperature,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(p.serverURL+"/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Text, nil
}

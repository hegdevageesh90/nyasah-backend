package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HuggingFaceProvider struct {
	apiKey string
	model  string
}

func NewHuggingFaceProvider(apiKey string, model string) *HuggingFaceProvider {
	return &HuggingFaceProvider{
		apiKey: apiKey,
		model:  model,
	}
}

func (p *HuggingFaceProvider) ProcessQuery(query string) (string, error) {
	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", p.model)

	payload := map[string]string{
		"inputs": query,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result) > 0 {
		return result[0]["generated_text"].(string), nil
	}

	return "", fmt.Errorf("no response generated")
}

func (p *HuggingFaceProvider) AnalyzeSentiment(text string) (float64, error) {
	url := "https://api-inference.huggingface.co/models/finiteautomata/bertweet-base-sentiment-analysis"

	payload := map[string]string{
		"inputs": text,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result [][]map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if len(result) > 0 && len(result[0]) > 0 {
		sentiment := result[0][0]["label"].(string)
		score := result[0][0]["score"].(float64)

		switch sentiment {
		case "POS":
			return score, nil
		case "NEG":
			return -score, nil
		default:
			return 0, nil
		}
	}

	return 0, fmt.Errorf("no sentiment analysis result")
}

func (p *HuggingFaceProvider) GenerateText(prompt string, maxTokens int, temperature float64) (string, error) {
	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", p.model)

	payload := map[string]interface{}{
		"inputs": prompt,
		"parameters": map[string]interface{}{
			"max_length":  maxTokens,
			"temperature": temperature,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result) > 0 {
		return result[0]["generated_text"].(string), nil
	}

	return "", fmt.Errorf("no response generated")
}

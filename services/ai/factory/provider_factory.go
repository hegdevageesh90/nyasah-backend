package factory

import (
	"fmt"
	"nyasah-backend/services/ai/providers"
)

type ProviderType string

const (
	OpenAI ProviderType = "openai"
	// TODO : bring back Claude
	//Claude      ProviderType = "claude"
	HuggingFace ProviderType = "huggingface"
	Llama       ProviderType = "llama"
)

func CreateProvider(providerType ProviderType, config map[string]string) (providers.Provider, error) {
	switch providerType {
	case OpenAI:
		apiKey, ok := config["api_key"]
		if !ok {
			return nil, fmt.Errorf("OpenAI API key not provided")
		}
		return providers.NewOpenAIProvider(apiKey), nil

	/*case Claude:
	apiKey, ok := config["api_key"]
	if !ok {
		return nil, fmt.Errorf("Claude API key not provided")
	}
	return providers.NewClaudeProvider(apiKey), nil*/

	case HuggingFace:
		apiKey, ok := config["api_key"]
		if !ok {
			return nil, fmt.Errorf("HuggingFace API key not provided")
		}
		model, ok := config["model"]
		if !ok {
			model = "gpt2" // Default model
		}
		return providers.NewHuggingFaceProvider(apiKey, model), nil

	case Llama:
		serverURL, ok := config["server_url"]
		if !ok {
			return nil, fmt.Errorf("Llama server URL not provided")
		}
		return providers.NewLlamaProvider(serverURL), nil

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

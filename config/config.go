package config

import (
	"fmt"
	"log"
	"nyasah-backend/services/ai/factory"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
	Provider    factory.ProviderType
	Model       string
	Temperature float64
	MaxTokens   int
}

func Load() (*Config, error) {
	err := godotenv.Load() // Load .env file if exists
	if err != nil {
		log.Printf("Warning: No .env file found. Using environment variables or defaults.\n")
	}

	temperature, err := getEnvAsFloat64("TEMPERATURE", 0.45) // Use realistic default for temperature
	if err != nil {
		return nil, err
	}

	maxTokens, err := getEnvAsInt("MAX_TOKENS", 1000) // Default max tokens
	if err != nil {
		return nil, err
	}

	provider, err := getProvider("PROVIDER", "meta") // Handle provider type
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		DatabaseURL: getEnv("DATABASE_URL", "nyasah.db"),
		Provider:    provider,
		Model:       getEnv("MODEL", "llama3.2"),
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) (int, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func getEnvAsFloat64(key string, fallback float64) (float64, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback, nil
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func getProvider(key, fallback string) (factory.ProviderType, error) {
	value := getEnv(key, fallback)
	switch value {
	case "meta":
		return factory.Llama, nil
	case "openai":
		return factory.OpenAI, nil
	case "huggingface":
		return factory.HuggingFace, nil
	default:
		return "", fmt.Errorf("invalid provider type: %s", value)
	}
}

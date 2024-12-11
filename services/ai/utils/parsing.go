package utils

import (
	"strings"
)

func ParseKeywords(text string) []string {
	keywords := strings.Split(text, ",")
	result := make([]string, 0)

	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword != "" {
			result = append(result, keyword)
		}
	}

	return result
}

func CalculateAverageScore(scores []float64) float64 {
	if len(scores) == 0 {
		return 0
	}

	var sum float64
	for _, score := range scores {
		sum += score
	}
	return sum / float64(len(scores))
}

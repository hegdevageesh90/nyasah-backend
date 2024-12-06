package utils

import (
	"nyasah-backend/models"
	"sort"
)

type Pattern struct {
	Type         string
	SampleSize   int
	SuccessRate  float64
	Distribution map[string]int
}

func AnalyzeContentPattern(proofs []models.SocialProof) Pattern {
	pattern := Pattern{
		Type:         "content",
		SampleSize:   len(proofs),
		Distribution: make(map[string]int),
	}

	successCount := 0
	for _, proof := range proofs {
		pattern.Distribution[proof.MediaType]++
		if proof.Performance.EngagementRate > 0.5 {
			successCount++
		}
	}

	if pattern.SampleSize > 0 {
		pattern.SuccessRate = float64(successCount) / float64(pattern.SampleSize)
	}

	return pattern
}

func AnalyzeTimingPattern(proofs []models.SocialProof) Pattern {
	pattern := Pattern{
		Type:         "timing",
		SampleSize:   len(proofs),
		Distribution: make(map[string]int),
	}

	successCount := 0
	for _, proof := range proofs {
		hour := proof.CreatedAt.Format("15")
		pattern.Distribution[hour]++
		if proof.Performance.EngagementRate > 0.5 {
			successCount++
		}
	}

	if pattern.SampleSize > 0 {
		pattern.SuccessRate = float64(successCount) / float64(pattern.SampleSize)
	}

	return pattern
}

func AnalyzePlacementPattern(proofs []models.SocialProof) Pattern {
	pattern := Pattern{
		Type:         "placement",
		SampleSize:   len(proofs),
		Distribution: make(map[string]int),
	}

	successCount := 0
	for _, proof := range proofs {
		pattern.Distribution[proof.Type]++
		if proof.Performance.EngagementRate > 0.5 {
			successCount++
		}
	}

	if pattern.SampleSize > 0 {
		pattern.SuccessRate = float64(successCount) / float64(pattern.SampleSize)
	}

	return pattern
}

func (p Pattern) GenerateRecommendation() string {
	switch p.Type {
	case "content":
		return p.generateContentRecommendation()
	case "timing":
		return p.generateTimingRecommendation()
	case "placement":
		return p.generatePlacementRecommendation()
	default:
		return "Optimize your social proof strategy based on performance data"
	}
}

func (p Pattern) generateContentRecommendation() string {
	// Find the most successful content type
	var maxCount int
	var bestType string
	for mediaType, count := range p.Distribution {
		if count > maxCount {
			maxCount = count
			bestType = mediaType
		}
	}

	switch bestType {
	case "image":
		return "Increase usage of visual content, particularly product images and customer photos"
	case "video":
		return "Focus on video content, such as customer testimonials and product demonstrations"
	default:
		return "Diversify your content mix with both visual and textual social proof"
	}
}

func (p Pattern) generateTimingRecommendation() string {
	// Find the most engaging time period
	type timeCount struct {
		hour  string
		count int
	}

	var times []timeCount
	for hour, count := range p.Distribution {
		times = append(times, timeCount{hour, count})
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i].count > times[j].count
	})

	if len(times) > 0 {
		return "Schedule social proof displays during peak engagement hours around " + times[0].hour + ":00"
	}

	return "Distribute social proof displays throughout the day"
}

func (p Pattern) generatePlacementRecommendation() string {
	// Find the most effective placement type
	var maxCount int
	var bestType string
	for proofType, count := range p.Distribution {
		if count > maxCount {
			maxCount = count
			bestType = proofType
		}
	}

	switch bestType {
	case "purchase":
		return "Emphasize recent purchase notifications to create urgency"
	case "review":
		return "Highlight customer reviews prominently on product pages"
	case "view":
		return "Display real-time viewer counts to show product popularity"
	default:
		return "Use a mix of social proof types across your site"
	}
}

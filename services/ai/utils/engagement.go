package utils

import (
	"nyasah-backend/models"
	"time"
)

func CalculateReviewEngagement(review models.Review) float64 {
	if review.Engagement.Views == 0 {
		return 0
	}

	engagementActions := float64(review.Engagement.Likes + review.Engagement.Shares)
	return engagementActions / float64(review.Engagement.Views)
}

func CalculateProofEngagement(proof models.SocialProof) float64 {
	if proof.Performance.Views == 0 {
		return 0
	}

	return float64(proof.Performance.Conversions) / float64(proof.Performance.Views)
}

func CalculateAverageSentiment(reviews []models.Review) float64 {
	if len(reviews) == 0 {
		return 0
	}

	var total float64
	for _, review := range reviews {
		total += review.Sentiment
	}
	return total / float64(len(reviews))
}

func CalculateAverageEngagement(reviews []models.Review, proofs []models.SocialProof) float64 {
	var total float64
	count := 0

	for _, review := range reviews {
		total += CalculateReviewEngagement(review)
		count++
	}

	for _, proof := range proofs {
		total += CalculateProofEngagement(proof)
		count++
	}

	if count == 0 {
		return 0
	}

	return total / float64(count)
}

func GroupByTimeFrames(reviews []models.Review, proofs []models.SocialProof) []TimeFrame {
	// Group data into weekly frames for the last 12 weeks
	frames := make([]TimeFrame, 12)
	now := time.Now()

	for i := range frames {
		frames[i].End = now.AddDate(0, 0, -7*i)
		frames[i].Start = frames[i].End.AddDate(0, 0, -7)
	}

	// Distribute reviews and proofs into frames
	for _, review := range reviews {
		for i, frame := range frames {
			if review.CreatedAt.After(frame.Start) && review.CreatedAt.Before(frame.End) {
				frames[i].Reviews = append(frames[i].Reviews, review)
				break
			}
		}
	}

	for _, proof := range proofs {
		for i, frame := range frames {
			if proof.CreatedAt.After(frame.Start) && proof.CreatedAt.Before(frame.End) {
				frames[i].Proofs = append(frames[i].Proofs, proof)
				break
			}
		}
	}

	return frames
}

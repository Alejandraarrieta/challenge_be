package mapper

import (
	domain "challenge_be/internal/domain/tweet"

	types "challenge_be/pkg/types/tweet/options"
)

func ToTweetResponse(tweet domain.Tweet) types.TweetResponse {
	var createdAt string
	if !tweet.CreatedAt.IsZero() {
		createdAt = tweet.CreatedAt.Format("02/01/2006 15:04")
	}

	return types.TweetResponse{
		ID:        tweet.ID,
		UserID:    tweet.UserID,
		Content:   tweet.Content,
		CreatedAt: createdAt,
	}
}

func MapTweetsToResponse(tweets []domain.Tweet) []types.TweetResponse {
	res := make([]types.TweetResponse, 0, len(tweets))
	for _, t := range tweets {
		res = append(res, ToTweetResponse(t))
	}
	return res
}

package tweet

import (
	domain_follow "challenge_be/internal/domain/follow"
	domain_tweet "challenge_be/internal/domain/tweet"
	"context"
)

type (
	GetTimelineUseCase interface {
		Execute(ctx context.Context, userID uint64) ([]domain_tweet.Tweet, error)
	}
	getTimelineUseCase struct {
		repository      domain_tweet.Repository
		cacheRepository domain_tweet.CacheRepository
		cacheFollowers  domain_follow.CacheRepository
	}
)

func NewGetTimelineUseCase(repository domain_tweet.Repository, cacheReposirory domain_tweet.CacheRepository, cacheFollowers domain_follow.CacheRepository) GetTimelineUseCase {
	return &getTimelineUseCase{repository, cacheReposirory, cacheFollowers}
}

func (g *getTimelineUseCase) Execute(ctx context.Context, userID uint64) ([]domain_tweet.Tweet, error) {

	followedIDs, err := g.cacheFollowers.GetFollowedUserIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(followedIDs) == 0 {
		return []domain_tweet.Tweet{}, nil // no sigue a nadie, no hay timeline
	}

	if len(followedIDs) > 0 {
		tweets, err := g.repository.GetTweetsByIDs(ctx, followedIDs)
		if err != nil {
			return nil, err
		}
		return tweets, nil
	}
	tweets, err := g.repository.ListTweetsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

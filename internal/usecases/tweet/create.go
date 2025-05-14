package tweet

import (
	domain "challenge_be/internal/domain/tweet"
	types "challenge_be/pkg/types/tweet/options"
	"context"
)

//go:generate mockgen -source=./create.go -destination=./mocks/create_mock.go -package=mocks

type (
	CreateUseCase interface {
		Execute(context.Context, types.InputCreateTweet) error
	}
	createUseCase struct {
		repository       domain.Repository
		cache_repository domain.CacheRepository
	}
)

func NewCreateUseCase(repository domain.Repository, cacheReposirory domain.CacheRepository) CreateUseCase {
	return &createUseCase{repository, cacheReposirory}
}

func (c *createUseCase) Execute(ctx context.Context, input types.InputCreateTweet) error {
	tweet := domain.Tweet{
		UserID:  input.UserID,
		Content: input.Content,
	}
	id, err := c.repository.Create(ctx, tweet)
	if err != nil {
		return err
	}
	tweet.ID = uint64(id) // Asignar el ID al tweet
	return c.cache_repository.PushTweetToTimeline(ctx, tweet)
}

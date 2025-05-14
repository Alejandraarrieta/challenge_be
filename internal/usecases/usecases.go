package usecases

import (
	repositories "challenge_be/internal/adapters/datasources/repositories"
	"challenge_be/internal/usecases/follow"
	"challenge_be/internal/usecases/tweet"
)

type UseCases struct {
	CreateTweetUseCase  tweet.CreateUseCase
	CreateFollowUseCase follow.CreateUseCase
	GetTimelineUseCase  tweet.GetTimelineUseCase
}

func CreateUsescases(
	repositories *repositories.Repositories,
) *UseCases {
	createTweetUsecase := tweet.NewCreateUseCase(
		repositories.TweetRepository, repositories.TweetCacheRepository)
	createFollowUsecase := follow.NewCreateUseCase(
		repositories.FollowRepository, repositories.FollowCacheRepository)
	getTimelineUsecase := tweet.NewGetTimelineUseCase(
		repositories.TweetRepository, repositories.TweetCacheRepository, repositories.FollowCacheRepository)
	return &UseCases{
		CreateTweetUseCase:  createTweetUsecase,
		CreateFollowUseCase: createFollowUsecase,
		GetTimelineUseCase:  getTimelineUsecase,
	}
}

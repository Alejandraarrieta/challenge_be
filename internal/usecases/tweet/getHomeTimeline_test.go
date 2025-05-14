package tweet_test

import (
	"context"
	"errors"
	"testing"

	mock_follow "challenge_be/internal/domain/follow/mocks"
	domain_tweet "challenge_be/internal/domain/tweet"
	mock_tweet "challenge_be/internal/domain/tweet/mocks"
	"challenge_be/internal/usecases/tweet"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetTimelineUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_tweet.NewMockRepository(ctrl)
	mockCache := mock_tweet.NewMockCacheRepository(ctrl)
	mockCacheFollowers := mock_follow.NewMockCacheRepository(ctrl)

	useCase := tweet.NewGetTimelineUseCase(mockRepo, mockCache, mockCacheFollowers)

	ctx := context.Background()
	userID := uint64(1)

	t.Run("success - with followed users", func(t *testing.T) {
		followedIDs := []uint64{2, 3}
		expectedTweets := []domain_tweet.Tweet{
			{ID: 1, UserID: 2, Content: "tweet 1"},
			{ID: 2, UserID: 3, Content: "tweet 2"},
		}

		mockCacheFollowers.EXPECT().GetFollowedUserIDs(ctx, userID).Return(followedIDs, nil)
		mockRepo.EXPECT().GetTweetsByIDs(ctx, followedIDs).Return(expectedTweets, nil)

		tweets, err := useCase.Execute(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTweets, tweets)
	})

	t.Run("success - no followed users", func(t *testing.T) {
		mockCacheFollowers.EXPECT().GetFollowedUserIDs(ctx, userID).Return([]uint64{}, nil)

		tweets, err := useCase.Execute(ctx, userID)
		assert.NoError(t, err)
		assert.Empty(t, tweets)
	})

	t.Run("error - cacheFollowers fails", func(t *testing.T) {
		mockCacheFollowers.EXPECT().GetFollowedUserIDs(ctx, userID).Return(nil, errors.New("cache error"))

		tweets, err := useCase.Execute(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, tweets)
		assert.Equal(t, "cache error", err.Error())
	})

	t.Run("error - repository fails", func(t *testing.T) {
		followedIDs := []uint64{2, 3}
		mockCacheFollowers.EXPECT().GetFollowedUserIDs(ctx, userID).Return(followedIDs, nil)
		mockRepo.EXPECT().GetTweetsByIDs(ctx, followedIDs).Return(nil, errors.New("repository error"))

		tweets, err := useCase.Execute(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, tweets)
		assert.Equal(t, "repository error", err.Error())
	})
}

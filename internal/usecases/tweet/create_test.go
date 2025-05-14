package tweet_test

import (
	domain "challenge_be/internal/domain/tweet"
	"challenge_be/internal/usecases/tweet"
	types "challenge_be/pkg/types/tweet/options"
	"context"
	"errors"
	"testing"

	"challenge_be/internal/domain/tweet/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockRepository(ctrl)
	mockCacheRepository := mocks.NewMockCacheRepository(ctrl)
	useCase := tweet.NewCreateUseCase(mockRepository, mockCacheRepository)

	ctx := context.Background()
	input := types.InputCreateTweet{
		UserID:  123,
		Content: "This is a test tweet",
	}

	t.Run("success", func(t *testing.T) {
		expectedTweet := domain.Tweet{
			UserID:  input.UserID,
			Content: input.Content,
		}

		mockRepository.EXPECT().Create(ctx, expectedTweet).Return(uint64(1), nil)
		expectedTweet.ID = 1
		mockCacheRepository.EXPECT().PushTweetToTimeline(ctx, expectedTweet).Return(nil)

		err := useCase.Execute(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedTweet := domain.Tweet{
			UserID:  input.UserID,
			Content: input.Content,
		}
		mockRepository.EXPECT().Create(ctx, expectedTweet).Return(uint64(0), errors.New("repository error"))
		mockCacheRepository.EXPECT().PushTweetToTimeline(gomock.Any(), gomock.Any()).Times(0)

		err := useCase.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
	})
	t.Run("cache repository error", func(t *testing.T) {
		expectedTweet := domain.Tweet{
			UserID:  input.UserID,
			Content: input.Content,
		}
		mockRepository.EXPECT().Create(ctx, expectedTweet).Return(uint64(0), nil)
		mockCacheRepository.EXPECT().PushTweetToTimeline(ctx, expectedTweet).Return(errors.New("cache error"))

		err := useCase.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "cache error", err.Error())
	})
}

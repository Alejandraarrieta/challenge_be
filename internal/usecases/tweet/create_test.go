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
		mockRepository.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, tweet *domain.Tweet) (uint64, error) {
				assert.Equal(t, input.UserID, tweet.UserID)
				assert.Equal(t, input.Content, tweet.Content)
				return 1, nil
			})

		mockCacheRepository.EXPECT().
			PushTweetToTimeline(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, tweet *domain.Tweet) error {
				assert.Equal(t, input.UserID, tweet.UserID)
				assert.Equal(t, input.Content, tweet.Content)
				return nil
			})

		err := useCase.Execute(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepository.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, tweet *domain.Tweet) (uint64, error) {
				assert.Equal(t, input.UserID, tweet.UserID)
				assert.Equal(t, input.Content, tweet.Content)
				return 0, errors.New("repository error")
			})

		// No se llama a PushTweetToTimeline si falla la creación
		mockCacheRepository.EXPECT().
			PushTweetToTimeline(gomock.Any(), gomock.Any()).
			Times(0)

		err := useCase.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
	})

	t.Run("cache repository error", func(t *testing.T) {
		mockRepository.EXPECT().
			Create(ctx, gomock.Any()).
			Return(uint64(1), nil)

		mockCacheRepository.EXPECT().
			PushTweetToTimeline(ctx, gomock.Any()).
			Return(errors.New("cache error"))

		err := useCase.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "cache error", err.Error())
	})
}

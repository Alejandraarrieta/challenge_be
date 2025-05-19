package follow_test

import (
	"context"
	"errors"
	"testing"

	domain "challenge_be/internal/domain/follow"
	"challenge_be/internal/domain/follow/mocks"
	"challenge_be/internal/usecases/follow"
	types "challenge_be/pkg/types/follow/options"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockRepository(ctrl)
	mockCacheRepository := mocks.NewMockCacheRepository(ctrl)
	useCase := follow.NewCreateUseCase(mockRepository, mockCacheRepository)

	ctx := context.Background()
	input := types.InputCreateFollow{
		FollowerID: 1,
		FolloweeID: 2,
	}

	t.Run("success", func(t *testing.T) {
		expectedFollow := &domain.Follow{
			FollowerID: input.FollowerID,
			FolloweeID: input.FolloweeID,
		}

		mockRepository.EXPECT().Create(ctx, expectedFollow).Return(nil)
		mockCacheRepository.EXPECT().AddFollowing(ctx, expectedFollow).Return(nil)

		err := useCase.Execute(ctx, input)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		expectedFollow := &domain.Follow{
			FollowerID: input.FollowerID,
			FolloweeID: input.FolloweeID,
		}

		mockRepository.EXPECT().Create(ctx, expectedFollow).Return(errors.New("repository error"))
		mockCacheRepository.EXPECT().AddFollowing(gomock.Any(), gomock.Any()).Times(0)

		err := useCase.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
	})

	t.Run("cache repository error", func(t *testing.T) {
		expectedFollow := &domain.Follow{
			FollowerID: input.FollowerID,
			FolloweeID: input.FolloweeID,
		}

		mockRepository.EXPECT().Create(ctx, expectedFollow).Return(nil)
		mockCacheRepository.EXPECT().AddFollowing(ctx, expectedFollow).Return(errors.New("cache error"))

		err := useCase.Execute(ctx, input)
		assert.Error(t, err)
		assert.Equal(t, "cache error", err.Error())
	})
}

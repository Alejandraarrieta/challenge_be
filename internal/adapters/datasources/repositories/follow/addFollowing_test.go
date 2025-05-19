package follow_test

import (
	"context"
	"errors"
	"testing"

	"challenge_be/internal/adapters/datasources/repositories/follow"
	domain "challenge_be/internal/domain/follow"
	"challenge_be/internal/platform/cache/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/redis/go-redis/v9"
)

func TestAddFollowing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := follow.NewCacheRepository(mockClient)
	ctx := context.Background()

	tests := map[string]struct {
		inputFollow *domain.Follow
		prepareMock func()
		expectErr   error
	}{
		"when followee is successfully added to the following set": {
			inputFollow: &domain.Follow{
				FollowerID: 100,
				FolloweeID: 200,
			},
			prepareMock: func() {
				key := "user:100:following"
				mockClient.EXPECT().
					SAdd(ctx, key, int64(200)).
					Return(redis.NewIntResult(int64(1), nil))
			},
			expectErr: nil,
		},
		"when Redis SAdd returns an error": {
			inputFollow: &domain.Follow{
				FollowerID: 300,
				FolloweeID: 400,
			},
			prepareMock: func() {
				key := "user:300:following"
				mockClient.EXPECT().
					SAdd(ctx, key, int64(400)).
					Return(redis.NewIntResult(0, errors.New("redis error")))
			},
			expectErr: errors.New("redis error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.prepareMock()
			err := repo.AddFollowing(ctx, tc.inputFollow)

			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

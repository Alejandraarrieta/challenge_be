package follow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"challenge_be/internal/adapters/datasources/repositories/follow"
	"challenge_be/internal/platform/cache/mocks"

	"github.com/redis/go-redis/v9"
)

func TestGetFollowedUserIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := follow.NewCacheRepository(mockClient)
	ctx := context.Background()
	userID := uint64(123)
	cacheKey := "user:123:following"

	tests := map[string]struct {
		prepareMock   func()
		expectErr     error
		expectUserIDs []uint64
	}{
		"when user IDs are successfully retrieved and parsed": {
			prepareMock: func() {
				mockClient.EXPECT().SMembers(ctx, cacheKey).Return(redis.NewStringSliceResult([]string{"1", "2", "3"}, nil))
			},
			expectUserIDs: []uint64{1, 2, 3},
		},
		"when no user IDs are found": {
			prepareMock: func() {
				mockClient.EXPECT().SMembers(ctx, cacheKey).Return(redis.NewStringSliceResult([]string{}, nil))
			},
			expectUserIDs: nil,
		},
		"when Redis SMembers returns an error": {
			prepareMock: func() {
				mockClient.EXPECT().SMembers(ctx, cacheKey).Return(redis.NewStringSliceResult(nil, errors.New("redis error")))
			},
			expectErr: errors.New("redis error"),
		},
		"when a cached member is not a valid uint64": {
			prepareMock: func() {
				mockClient.EXPECT().SMembers(ctx, cacheKey).Return(redis.NewStringSliceResult([]string{"1", "invalid", "3"}, nil))
			},
			expectUserIDs: []uint64{1, 3},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.prepareMock()
			ids, err := repo.GetFollowedUserIDs(ctx, userID)

			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
				assert.Nil(t, ids)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectUserIDs, ids)
			}
		})
	}
}

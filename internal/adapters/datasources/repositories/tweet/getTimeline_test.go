package tweet_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"challenge_be/internal/adapters/datasources/repositories/tweet"
	"challenge_be/internal/platform/cache/mocks"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := tweet.NewCacheRepository(mockClient)
	ctx := context.Background()
	userID := uint64(123)
	cacheKey := fmt.Sprintf("user:%d:timeline", userID)

	tests := map[string]struct {
		prepareMock func()
		expectErr   error
		expectIDs   []uint64
	}{
		"when timeline IDs are successfully retrieved and parsed": {
			prepareMock: func() {
				mockClient.EXPECT().LRange(ctx, cacheKey, int64(0), int64(-1)).Return(redis.NewStringSliceResult([]string{"1", "2", "3"}, nil))
			},
			expectIDs: []uint64{1, 2, 3},
		},
		"when no timeline IDs are in cache": {
			prepareMock: func() {
				mockClient.EXPECT().LRange(ctx, cacheKey, int64(0), int64(-1)).Return(redis.NewStringSliceResult([]string{}, nil))
			},
			expectIDs: nil,
		},
		"when Redis LRange returns an error": {
			prepareMock: func() {
				mockClient.EXPECT().LRange(ctx, cacheKey, int64(0), int64(-1)).Return(redis.NewStringSliceResult(nil, errors.New("redis error")))
			},
			expectErr: errors.New("redis error"),
		},
		"when a cached ID is not a valid uint64": {
			prepareMock: func() {
				mockClient.EXPECT().LRange(ctx, cacheKey, int64(0), int64(-1)).Return(redis.NewStringSliceResult([]string{"1", "invalid", "3"}, nil))
			},
			expectErr: fmt.Errorf("error parsing tweet ID \"invalid\": strconv.ParseUint: parsing \"invalid\": invalid syntax"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.prepareMock()
			ids, err := repo.GetTimeline(ctx, userID)

			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
				assert.Nil(t, ids)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectIDs, ids)
			}
		})
	}
}

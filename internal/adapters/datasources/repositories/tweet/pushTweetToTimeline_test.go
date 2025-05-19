package tweet_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"challenge_be/internal/adapters/datasources/repositories/tweet"
	domain "challenge_be/internal/domain/tweet"
	"challenge_be/internal/platform/cache/mocks"

	"github.com/redis/go-redis/v9"
)

func TestPushTweetToTimeline(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockRedisClientInterface(ctrl)
	repo := tweet.NewCacheRepository(mockClient)
	ctx := context.Background()
	testTweet := &domain.Tweet{
		ID:     999,
		UserID: 555,
	}
	cacheKey := fmt.Sprintf("user:%d:timeline", testTweet.UserID)
	tweetIDStr := strconv.FormatUint(testTweet.ID, 10)

	tests := map[string]struct {
		prepareMock func()
		expectErr   error
	}{
		"when tweet ID is successfully pushed to timeline": {
			prepareMock: func() {
				mockClient.EXPECT().LPush(ctx, cacheKey, tweetIDStr).Return(redis.NewIntResult(int64(1), nil))
			},
			expectErr: nil,
		},
		"when Redis LPush returns an error": {
			prepareMock: func() {
				mockClient.EXPECT().LPush(ctx, cacheKey, tweetIDStr).Return(redis.NewIntResult(int64(0), errors.New("redis error")))
			},
			expectErr: errors.New("redis error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.prepareMock()
			err := repo.PushTweetToTimeline(ctx, testTweet)

			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

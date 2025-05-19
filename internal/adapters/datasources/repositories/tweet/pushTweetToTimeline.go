package tweet

import (
	domain "challenge_be/internal/domain/tweet"
	"context"
	"fmt"
)

func (r *cacheRepository) PushTweetToTimeline(ctx context.Context, tweet *domain.Tweet) error {
	key := fmt.Sprintf("user:%d:timeline", tweet.UserID)
	return r.client.LPush(ctx, key, fmt.Sprintf("%d", tweet.ID)).Err()
}

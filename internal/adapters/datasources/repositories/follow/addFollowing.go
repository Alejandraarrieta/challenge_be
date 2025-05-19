package follow

import (
	domain "challenge_be/internal/domain/follow"
	"context"
	"fmt"
)

func (r *cacheRepository) AddFollowing(ctx context.Context, follow *domain.Follow) error {
	key := fmt.Sprintf("user:%d:following", follow.FollowerID)
	fmt.Printf("CACHE: Attempting to SAdd key: %s, member: %d\n", key, follow.FolloweeID)
	err := r.client.SAdd(ctx, key, int64(follow.FolloweeID)).Err()
	if err != nil {
		fmt.Printf("CACHE: Error during SAdd: %v\n", err)
		return err
	}
	fmt.Println("CACHE: SAdd successful")
	return nil
}

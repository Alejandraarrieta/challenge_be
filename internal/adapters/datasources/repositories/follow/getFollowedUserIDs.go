package follow

import (
	"context"
	"fmt"
	"strconv"
)

func (r *cacheRepository) GetFollowedUserIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	key := fmt.Sprintf("user:%d:following", userID)
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var ids []uint64
	for _, m := range members {
		id, err := strconv.ParseUint(m, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

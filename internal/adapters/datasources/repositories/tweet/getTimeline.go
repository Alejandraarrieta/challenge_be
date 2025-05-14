package tweet

import (
	"context"
	"fmt"
	"strconv"
)

func (r *cacheRepository) GetTimeline(ctx context.Context, userID uint64) ([]uint64, error) {
	key := fmt.Sprintf("user:%d:timeline", userID)
	idStrings, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var ids []uint64
	for _, idStr := range idStrings {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing tweet ID %q: %w", idStr, err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

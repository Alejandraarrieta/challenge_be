package tweet

import (
	domain "challenge_be/internal/domain/tweet"
	"challenge_be/internal/platform/cache"
)

type cacheRepository struct {
	client cache.RedisClientInterface
}

func NewCacheRepository(client cache.RedisClientInterface) domain.CacheRepository {
	return &cacheRepository{client}
}

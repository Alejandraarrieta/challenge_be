package repositories

import (
	"challenge_be/internal/adapters/datasources/repositories/follow"
	"challenge_be/internal/adapters/datasources/repositories/tweet"
	domain_follow "challenge_be/internal/domain/follow"
	domain_tweet "challenge_be/internal/domain/tweet"
	"challenge_be/internal/platform/cache"
	"database/sql"
)

type Repositories struct {
	TweetRepository       domain_tweet.Repository
	TweetCacheRepository  domain_tweet.CacheRepository
	FollowRepository      domain_follow.Repository
	FollowCacheRepository domain_follow.CacheRepository
}

func CreateRepository(db *sql.DB, redisClient cache.RedisClientInterface) *Repositories {

	tweetRepository := tweet.NewRepository(db)
	followRepository := follow.NewRepository(db)
	followCacheRepository := follow.NewCacheRepository(redisClient)
	tweetCacheRepository := tweet.NewCacheRepository(redisClient)
	return &Repositories{
		TweetRepository:       tweetRepository,
		TweetCacheRepository:  tweetCacheRepository,
		FollowRepository:      followRepository,
		FollowCacheRepository: followCacheRepository,
	}
}

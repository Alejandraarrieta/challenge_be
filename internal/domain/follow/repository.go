package domain

import (
	"context"
)

type Repository interface {
	Create(context.Context, Follow) error
}

type CacheRepository interface {
	AddFollowing(context.Context, Follow) error
	GetFollowedUserIDs(context.Context, uint64) ([]uint64, error)
}

package domain

import (
	"context"
)

type Repository interface {
	Create(context.Context, *Tweet) (uint64, error)
	ListTweetsByUserID(context.Context, uint64) ([]Tweet, error)
	GetTweetsByIDs(context.Context, []uint64) ([]Tweet, error)
}

type CacheRepository interface {
	PushTweetToTimeline(context.Context, *Tweet) error
	GetTimeline(context.Context, uint64) ([]uint64, error)
}

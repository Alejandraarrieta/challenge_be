package follow

import (
	domain "challenge_be/internal/domain/follow"
	types "challenge_be/pkg/types/follow/options"
	"context"
)

//go:generate mockgen -source=./create.go -destination=./mocks/create_mock.go -package=mocks

type (
	CreateUseCase interface {
		Execute(context.Context, types.InputCreateFollow) error
	}
	createUseCase struct {
		repository       domain.Repository
		cache_repository domain.CacheRepository
	}
)

func NewCreateUseCase(repository domain.Repository, cacheReposirory domain.CacheRepository) CreateUseCase {
	return &createUseCase{repository, cacheReposirory}
}

func (c *createUseCase) Execute(ctx context.Context, input types.InputCreateFollow) error {
	follow := domain.Follow{
		FollowerID: input.FollowerID,
		FolloweeID: input.FolloweeID,
	}

	if err := c.repository.Create(ctx, follow); err != nil {
		return err
	}

	return c.cache_repository.AddFollowing(ctx, follow)
}

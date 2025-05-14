package repositories_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	repo "challenge_be/internal/adapters/datasources/repositories"
	"challenge_be/internal/adapters/datasources/repositories/follow"
	"challenge_be/internal/adapters/datasources/repositories/tweet"
	cache_mock "challenge_be/internal/platform/cache/mocks"
)

func TestCreateRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeCache := cache_mock.NewMockRedisClientInterface(ctrl)

	actual := repo.CreateRepository(db, fakeCache)

	// Verifica que los repositorios no sean nil y sean del tipo esperado
	assert.NotNil(t, actual.TweetRepository, "TweetRepository should not be nil")
	assert.IsType(t, tweet.NewRepository(db), actual.TweetRepository, "TweetRepository should be of the expected type")

	assert.NotNil(t, actual.FollowRepository, "FollowRepository should not be nil")
	assert.IsType(t, follow.NewRepository(db), actual.FollowRepository, "FollowRepository should be of the expected type")

	assert.NotNil(t, actual.FollowCacheRepository, "FollowCacheRepository should not be nil")
	assert.IsType(t, follow.NewCacheRepository(fakeCache), actual.FollowCacheRepository, "FollowCacheRepository should be of the expected type")
}

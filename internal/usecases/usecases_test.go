package usecases_test

/*import (
	"testing"

	"challenge_be/internal/adapters/datasources/repositories"
	tweet_repository "challenge_be/internal/domain/tweet/mocks"
	"challenge_be/internal/usecases"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateRepositories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Initialize the mock database
	repositories := &repositories.Repositories{
		TweetRepository: tweet_repository.NewMockRepository(t),
	}

	expext := usecases.UseCases{
		CreateTweetUseCase: tweet_repository.NewMockRepository(repositories.TweetRepository),
	}

	actual := usecases.CreateUsescases(repositories)

	assert.Equalf(t, expext, actual, "expected %v, got %v", expext, actual)

}*/

import (
	"challenge_be/internal/adapters/datasources/repositories"
	"challenge_be/internal/domain/tweet/mocks"
	"challenge_be/internal/usecases"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUseCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTweetRepo := mocks.NewMockRepository(ctrl)

	repos := &repositories.Repositories{
		TweetRepository: mockTweetRepo,
	}

	uc := usecases.CreateUsescases(repos)

	assert.NotNil(t, uc)
	assert.NotNil(t, uc.CreateTweetUseCase)
}

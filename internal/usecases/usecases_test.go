package usecases_test

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

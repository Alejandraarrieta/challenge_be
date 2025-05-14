package tweet_test

import (
	"challenge_be/internal/adapters/web/handlers/tweet"
	domain "challenge_be/internal/domain/tweet"
	"challenge_be/internal/usecases/tweet/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewGetTimelineHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 with tweets", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockGetTimelineUseCase(ctrl)
		handler := tweet.NewGetTimelineHandler(mockUseCase)

		expectedTweets := []domain.Tweet{
			{ID: 1, UserID: 123, Content: "Hola mundo"},
			{ID: 2, UserID: 124, Content: "Segundo tweet"},
		}

		mockUseCase.
			EXPECT().
			Execute(gomock.Any(), uint64(123)).
			Return(expectedTweets, nil)

		req := httptest.NewRequest(http.MethodGet, "/tweets/123", nil)
		w := httptest.NewRecorder()

		r := gin.New()
		r.GET("/tweets/:user_id", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Hola mundo")
	})

	t.Run("should return 400 when user_id is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tweets/", nil)
		w := httptest.NewRecorder()

		r := gin.New()
		r.GET("/tweets/", tweet.NewGetTimelineHandler(nil))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // porque el path no hace match
	})

	t.Run("should return 400 when user_id is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tweets/abc", nil)
		w := httptest.NewRecorder()

		r := gin.New()
		r.GET("/tweets/:user_id", tweet.NewGetTimelineHandler(nil))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user_id must be a number")
	})

	t.Run("should return 500 when usecase fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockGetTimelineUseCase(ctrl)
		handler := tweet.NewGetTimelineHandler(mockUseCase)

		mockUseCase.
			EXPECT().
			Execute(gomock.Any(), uint64(123)).
			Return(nil, errors.New("usecase error"))

		req := httptest.NewRequest(http.MethodGet, "/tweets/123", nil)
		w := httptest.NewRecorder()

		r := gin.New()
		r.GET("/tweets/:user_id", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to get timeline")
	})
}

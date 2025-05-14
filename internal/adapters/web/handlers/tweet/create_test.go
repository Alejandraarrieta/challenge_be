package tweet_test

import (
	"bytes"
	"challenge_be/internal/adapters/web/handlers/tweet"
	"challenge_be/internal/usecases/tweet/mocks"
	types "challenge_be/pkg/types/tweet/options"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewCreateTweetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 201 on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockCreateUseCase(ctrl)
		handler := tweet.NewCreateTweetHandler(mockUseCase)

		inputJSON := `{"user_id": 123, "content": "Hello, world!"}`

		req := httptest.NewRequest(http.MethodPost, "/tweets", bytes.NewBufferString(inputJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockUseCase.
			EXPECT().
			Execute(gomock.Any(), types.InputCreateTweet{UserID: 123, Content: "Hello, world!"}).
			Return(nil)

		r := gin.New()
		r.POST("/tweets", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Tweet created successfully")
	})

	t.Run("should return 400 on invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockCreateUseCase(ctrl)
		handler := tweet.NewCreateTweetHandler(mockUseCase)

		invalidJSON := `{"user_id": "not_a_number", "content": 123}`

		req := httptest.NewRequest(http.MethodPost, "/tweets", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := gin.New()
		r.POST("/tweets", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("should return 500 when usecase fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockCreateUseCase(ctrl)
		handler := tweet.NewCreateTweetHandler(mockUseCase)

		inputJSON := `{"user_id": 123, "content": "Failing tweet"}`

		req := httptest.NewRequest(http.MethodPost, "/tweets", bytes.NewBufferString(inputJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockUseCase.
			EXPECT().
			Execute(gomock.Any(), types.InputCreateTweet{UserID: 123, Content: "Failing tweet"}).
			Return(errors.New("something went wrong"))

		r := gin.New()
		r.POST("/tweets", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create tweet")
	})
}

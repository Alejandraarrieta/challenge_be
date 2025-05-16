package follow_test

import (
	"bytes"
	"challenge_be/internal/adapters/web/handlers/follow"
	"challenge_be/internal/usecases/follow/mocks"
	types "challenge_be/pkg/types/follow/options"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewCreateFollowHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 201 on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockCreateUseCase(ctrl)
		handler := follow.NewCreateFollowHandler(mockUseCase)

		inputJSON := `{"follower_id": 123, "followee_id": 456}`

		req := httptest.NewRequest(http.MethodPost, "/follows", bytes.NewBufferString(inputJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Espera que el método Execute sea llamado con los valores correctos
		mockUseCase.
			EXPECT().
			Execute(gomock.Any(), types.InputCreateFollow{FollowerID: 123, FolloweeID: 456}).
			Return(nil).Times(1) // Espera que se llame exactamente una vez

		r := gin.New()
		r.POST("/follows", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Follow created successfully")
	})

	t.Run("should return 400 on invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockCreateUseCase(ctrl)
		handler := follow.NewCreateFollowHandler(mockUseCase)

		// JSON inválido
		invalidJSON := `{"follower_id": "not_a_number", "followee_id": 456}`

		req := httptest.NewRequest(http.MethodPost, "/follows", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := gin.New()
		r.POST("/follows", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("should return 500 when usecase fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockCreateUseCase(ctrl)
		handler := follow.NewCreateFollowHandler(mockUseCase)

		inputJSON := `{"follower_id": 123, "followee_id": 456}`

		req := httptest.NewRequest(http.MethodPost, "/follows", bytes.NewBufferString(inputJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Simula un error del usecase
		mockUseCase.
			EXPECT().
			Execute(gomock.Any(), types.InputCreateFollow{FollowerID: 123, FolloweeID: 456}).
			Return(errors.New("something went wrong")).
			Times(1)

		r := gin.New()
		r.POST("/follows", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to create follow")
	})

	t.Run("should return 400 when required fields are missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUseCase := mocks.NewMockCreateUseCase(ctrl)
		handler := follow.NewCreateFollowHandler(mockUseCase)

		// Falta el campo "followee_id"
		inputJSON := `{"follower_id": 123}`

		req := httptest.NewRequest(http.MethodPost, "/follows", bytes.NewBufferString(inputJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r := gin.New()
		r.POST("/follows", handler)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"FolloweeID":"El campo es requerido"`)
	})
}

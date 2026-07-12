package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"jdgonzalez907/users-api/internal/infrastructure/controllers"
	mockApp "jdgonzalez907/users-api/mocks/application"

	"github.com/stretchr/testify/assert"
)

func TestRouterAndMiddleware(t *testing.T) {
	mockFindUseCase := mockApp.NewMockFindUserByIdUseCase(t)
	mockCreateUseCase := mockApp.NewMockCreateUserUseCase(t)

	findController := controllers.NewFindUserByIDController(mockFindUseCase)
	createController := controllers.NewCreateUserController(mockCreateUseCase)

	router := controllers.NewRouter(findController, createController)

	t.Run("JSONContentTypeMiddleware sets header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()

		// Set mock expectation for path /users/1
		mockFindUseCase.EXPECT().Execute(1).Return(nil, http.ErrNoLocation).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
}

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
	mockDeleteUseCase := mockApp.NewMockDeleteUserUseCase(t)
	mockUpdatePIUseCase := mockApp.NewMockUpdateUserPersonalInformationUseCase(t)

	findController := controllers.NewFindUserByIDController(mockFindUseCase)
	createController := controllers.NewCreateUserController(mockCreateUseCase)
	deleteController := controllers.NewDeleteUserController(mockDeleteUseCase)
	updatePIController := controllers.NewUpdateUserPersonalInformationController(mockUpdatePIUseCase)

	router := controllers.NewRouter(controllers.RouterParams{
		FindUserByID:              findController,
		CreateUser:                createController,
		DeleteUser:                deleteController,
		UpdatePersonalInformation: updatePIController,
	})

	t.Run("JSONContentTypeMiddleware sets header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()

		mockFindUseCase.EXPECT().Execute(1).Return(nil, http.ErrNoLocation).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
}

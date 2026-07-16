package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestRouterAndMiddleware(t *testing.T) {
	mockFindUseCase := mockApp.NewMockFindUserByIDUseCase(t)
	mockCreateUseCase := mockApp.NewMockCreateUserUseCase(t)
	mockDeleteUseCase := mockApp.NewMockDeleteUserUseCase(t)
	mockUpdatePIUseCase := mockApp.NewMockUpdateUserPersonalInformationUseCase(t)
	mockFindPaginatedUseCase := mockApp.NewMockFindUsersPaginatedUseCase(t)
	mockUpdateEmailUseCase := mockApp.NewMockChangeUserEmailUseCase(t)
	mockUpdatePhoneUseCase := mockApp.NewMockChangeUserPhoneUseCase(t)

	findController := controllers.NewFindUserByIDController(mockFindUseCase)
	createController := controllers.NewCreateUserController(mockCreateUseCase)
	deleteController := controllers.NewDeleteUserController(mockDeleteUseCase)
	updatePIController := controllers.NewUpdateUserPersonalInformationController(mockUpdatePIUseCase)
	findPaginatedController := controllers.NewFindUsersPaginatedController(mockFindPaginatedUseCase)
	updateEmailController := controllers.NewChangeUserEmailController(mockUpdateEmailUseCase)
	updatePhoneController := controllers.NewChangeUserPhoneController(mockUpdatePhoneUseCase)

	router := controllers.NewRouter(controllers.RouterParams{
		FindUserByID:              findController,
		CreateUser:                createController,
		DeleteUser:                deleteController,
		UpdatePersonalInformation: updatePIController,
		FindUsersPaginated:        findPaginatedController,
		ChangeEmail:               updateEmailController,
		ChangePhone:               updatePhoneController,
	})

	t.Run("JSONContentTypeMiddleware sets header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		req.Header.Set("Authorization", "1")
		rec := httptest.NewRecorder()

		mockFindUseCase.EXPECT().Execute(mock.Anything, int64(1)).Return(nil, http.ErrNoLocation).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
}

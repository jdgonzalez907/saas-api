package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"
	mockApp "jdgonzalez907/saas-api/mocks/application"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRouterAndMiddleware(t *testing.T) {
	mockFindUseCase := mockApp.NewMockFindUserByIDUseCase(t)
	mockCreateUseCase := mockApp.NewMockCreateUserUseCase(t)
	mockDeleteUseCase := mockApp.NewMockDeleteUserUseCase(t)
	mockUpdatePIUseCase := mockApp.NewMockUpdateUserPersonalInformationUseCase(t)
	mockFindPaginatedUseCase := mockApp.NewMockFindUsersPaginatedUseCase(t)
	mockUpdateEmailUseCase := mockApp.NewMockUpdateUserEmailUseCase(t)
	mockUpdatePhoneUseCase := mockApp.NewMockUpdateUserPhoneUseCase(t)

	findController := controllers.NewFindUserByIDController(mockFindUseCase)
	createController := controllers.NewCreateUserController(mockCreateUseCase)
	deleteController := controllers.NewDeleteUserController(mockDeleteUseCase)
	updatePIController := controllers.NewUpdateUserPersonalInformationController(mockUpdatePIUseCase)
	findPaginatedController := controllers.NewFindUsersPaginatedController(mockFindPaginatedUseCase)
	updateEmailController := controllers.NewUpdateUserEmailController(mockUpdateEmailUseCase)
	updatePhoneController := controllers.NewUpdateUserPhoneController(mockUpdatePhoneUseCase)

	router := controllers.NewRouter(controllers.RouterParams{
		FindUserByID:              findController,
		CreateUser:                createController,
		DeleteUser:                deleteController,
		UpdatePersonalInformation: updatePIController,
		FindUsersPaginated:        findPaginatedController,
		UpdateEmail:               updateEmailController,
		UpdatePhone:               updatePhoneController,
	})

	t.Run("JSONContentTypeMiddleware sets header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()

		mockFindUseCase.EXPECT().Execute(mock.Anything, int64(1)).Return(nil, http.ErrNoLocation).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
}

package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
	"jdgonzalez907/users-api/internal/infrastructure/controllers"
	mockApp "jdgonzalez907/users-api/mocks/application"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestFindUserByIDController_Handle(t *testing.T) {
	identification, _ := domain.NewIdentification(domain.IdType_CC, "123456789")
	phone, _ := domain.NewPhone("57", "3112223344")
	email, _ := domain.NewEmail("test@example.com")
	address, _ := domain.NewAddress("Street 1", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-20, 0, 0))

	validUser, err := domain.NewUser(domain.UserParams{
		ID:             1,
		Identification: identification,
		FirstName:      "John",
		LastName:       "Doe",
		Phone:          phone,
		Email:          &email,
		Address:        &address,
		BirthDate:      &birthDate,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	assert.NoError(t, err)

	testCases := []struct {
		testName       string
		routeParamID   string
		setupMock      func(m *mockApp.MockFindUserByIdUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - user found",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindUserByIdUseCase) {
				m.EXPECT().Execute(1).Return(validUser, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName:     "fail - route parameter is not an integer",
			routeParamID: "abc",
			setupMock:    func(m *mockApp.MockFindUserByIdUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:     "fail - route parameter is empty",
			routeParamID: "",
			setupMock:    func(m *mockApp.MockFindUserByIdUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id is missing",
		},
		{
			testName:     "fail - route parameter is negative",
			routeParamID: "-5",
			setupMock:    func(m *mockApp.MockFindUserByIdUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:     "fail - user not found",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindUserByIdUseCase) {
				m.EXPECT().Execute(1).Return(nil, domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrUserNotFound.Error(),
		},
		{
			testName:     "fail - internal server error from usecase",
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindUserByIdUseCase) {
				m.EXPECT().Execute(1).Return(nil, errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockFindUserByIdUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewFindUserByIDController(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.routeParamID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			controller.Handle(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				var jsonResponse map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &jsonResponse)
				assert.NoError(t, err)
				assert.Contains(t, jsonResponse["message"], tc.expectedBody)
			} else {
				var responseDTO domain.UserDTO
				err := json.Unmarshal(rec.Body.Bytes(), &responseDTO)
				assert.NoError(t, err)
				assert.Equal(t, validUser.ID(), responseDTO.ID)
			}
		})
	}
}

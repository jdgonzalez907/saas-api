package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
	"jdgonzalez907/saas-api/internal/users/domain"
	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestFindUserByIDController_Handle(t *testing.T) {
	identification, _ := domain.NewIdentification(domain.IDTypeCC, "123456789")
	phone, _ := domain.NewPhone("57", "3112223344")
	email, _ := domain.NewEmail("test@example.com")
	address, _ := domain.NewAddress("Street 1", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-20, 0, 0).Format("2006-01-02"))

	personalInfo, _ := domain.NewPersonalInformation(
		identification,
		"John",
		"Doe",
		&address,
		&birthDate,
	)

	validUser, err := domain.NewUser(
		1,
		personalInfo,
		phone,
		&email,
		time.Now(),
		time.Now(),
	)
	assert.NoError(t, err)

	testCases := []struct {
		testName       string
		authUserID     any
		routeParamID   string
		setupMock      func(m *mockApp.MockFindUserByIDUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:     "success - user found",
			authUserID:   int64(1),
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindUserByIDUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), int64(1)).Return(validUser, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName:       "fail - unauthenticated",
			authUserID:     nil,
			routeParamID:   "1",
			setupMock:      func(_ *mockApp.MockFindUserByIDUseCase) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   sharedHttp.ErrUnauthenticated.Error(),
		},
		{
			testName:       "fail - route parameter is not an integer",
			authUserID:     int64(1),
			routeParamID:   "abc",
			setupMock:      func(_ *mockApp.MockFindUserByIDUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:       "fail - route parameter is empty",
			authUserID:     int64(1),
			routeParamID:   "",
			setupMock:      func(_ *mockApp.MockFindUserByIDUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id is missing",
		},
		{
			testName:       "fail - route parameter is negative",
			authUserID:     int64(1),
			routeParamID:   "-5",
			setupMock:      func(_ *mockApp.MockFindUserByIDUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter id must be a positive integer",
		},
		{
			testName:     "fail - user not found",
			authUserID:   int64(1),
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindUserByIDUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), int64(1)).Return(nil, domain.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrUserNotFound.Error(),
		},
		{
			testName:     "fail - internal server error from usecase",
			authUserID:   int64(1),
			routeParamID: "1",
			setupMock: func(m *mockApp.MockFindUserByIDUseCase) {
				m.EXPECT().Execute(mock.Anything, int64(1), int64(1)).Return(nil, errors.New("database connection lost"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockFindUserByIDUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewFindUserByIDController(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			if tc.authUserID != nil {
				req.Header.Set("Authorization", strconv.FormatInt(tc.authUserID.(int64), 10))
			}
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.routeParamID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			handler := sharedHttp.Protected(controller.Handle)
			handler.ServeHTTP(rec, req)

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

package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
	"jdgonzalez907/saas-api/internal/users/domain"
	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"
	mockApp "jdgonzalez907/saas-api/mocks/application"
)

func TestFindUsersPaginatedController_Handle(t *testing.T) {
	identification, _ := domain.NewIdentification(domain.IDTypeCC, "123456")
	phone, _ := domain.NewPhone("57", "987654321")
	email, _ := domain.NewEmail("test@example.com")
	address, _ := domain.NewAddress("St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-25, 0, 0).Format("2006-01-02"))

	personalInfo, _ := domain.NewPersonalInformation(
		identification,
		"John",
		"Doe",
		&address,
		&birthDate,
	)

	u, err := domain.NewUser(
		int64(1),
		personalInfo,
		phone,
		&email,
		time.Now(),
		time.Now(),
	)
	assert.NoError(t, err)

	nextCursor := int64(1)
	paginatedUsers := domain.NewPaginatedUsers([]*domain.User{u}, &nextCursor)

	testCases := []struct {
		testName       string
		authUserID     any
		urlQuery       string
		setupMock      func(m *mockApp.MockFindUsersPaginatedUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName:   "success - default pagination without query params",
			authUserID: int64(1),
			urlQuery:   "",
			setupMock: func(m *mockApp.MockFindUsersPaginatedUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.MatchedBy(func(p domain.Pagination) bool {
					return p.LastID() == nil && p.Limit() == 10
				})).Return(paginatedUsers, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"users"`,
		},
		{
			testName:   "success - valid cursor and limit",
			authUserID: int64(1),
			urlQuery:   "?cursor=5&limit=25",
			setupMock: func(m *mockApp.MockFindUsersPaginatedUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.MatchedBy(func(p domain.Pagination) bool {
					return p.LastID() != nil && *p.LastID() == 5 && p.Limit() == 25
				})).Return(paginatedUsers, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"users"`,
		},
		{
			testName:       "fail - unauthenticated",
			authUserID:     nil,
			urlQuery:       "",
			setupMock:      func(_ *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   sharedHttp.ErrUnauthenticated.Error(),
		},
		{
			testName:       "fail - invalid cursor string",
			authUserID:     int64(1),
			urlQuery:       "?cursor=abc",
			setupMock:      func(_ *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter cursor must be a positive integer",
		},
		{
			testName:       "fail - invalid cursor negative",
			authUserID:     int64(1),
			urlQuery:       "?cursor=-5",
			setupMock:      func(_ *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter cursor must be a positive integer",
		},
		{
			testName:       "fail - invalid limit string",
			authUserID:     int64(1),
			urlQuery:       "?limit=abc",
			setupMock:      func(_ *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter limit must be a positive integer",
		},
		{
			testName:       "fail - invalid limit negative",
			authUserID:     int64(1),
			urlQuery:       "?limit=-10",
			setupMock:      func(_ *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter limit must be a positive integer",
		},
		{
			testName:       "fail - invalid limit value not allowed",
			authUserID:     int64(1),
			urlQuery:       "?limit=12",
			setupMock:      func(_ *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidPaginationLimit.Error(),
		},
		{
			testName:   "fail - usecase execution error",
			authUserID: int64(1),
			urlQuery:   "",
			setupMock: func(m *mockApp.MockFindUsersPaginatedUseCase) {
				m.EXPECT().Execute(mock.Anything, mock.Anything).Return(domain.PaginatedUsers{}, domain.ErrFindingUsers)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockUseCase := mockApp.NewMockFindUsersPaginatedUseCase(t)
			tc.setupMock(mockUseCase)

			controller := controllers.NewFindUsersPaginatedController(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/users"+tc.urlQuery, nil)
			if tc.authUserID != nil {
				req.Header.Set("Authorization", strconv.FormatInt(tc.authUserID.(int64), 10))
			}
			rec := httptest.NewRecorder()

			handler := sharedHttp.Protected(controller.Handle)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
			}
		})
	}
}

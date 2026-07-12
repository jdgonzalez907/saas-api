package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
	"jdgonzalez907/users-api/internal/infrastructure/controllers"
	mockApp "jdgonzalez907/users-api/mocks/application"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindUsersPaginatedController_Handle(t *testing.T) {
	identification, _ := domain.NewIdentification(domain.IdType_CC, "123456")
	phone, _ := domain.NewPhone("57", "987654321")
	email, _ := domain.NewEmail("test@example.com")
	address, _ := domain.NewAddress("St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-25, 0, 0))

	userParams := domain.UserParams{
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
	}
	u, err := domain.NewUser(userParams)
	assert.NoError(t, err)

	nextCursor := 1
	paginatedUsers := domain.NewPaginatedUsers([]*domain.User{u}, &nextCursor)

	testCases := []struct {
		testName       string
		urlQuery       string
		setupMock      func(m *mockApp.MockFindUsersPaginatedUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			testName: "success - default pagination without query params",
			urlQuery: "",
			setupMock: func(m *mockApp.MockFindUsersPaginatedUseCase) {
				m.EXPECT().Execute(mock.MatchedBy(func(p domain.Pagination) bool {
					return p.LastID() == nil && p.Limit() == 10
				})).Return(paginatedUsers, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"users"`,
		},
		{
			testName: "success - valid cursor and limit",
			urlQuery: "?cursor=5&limit=25",
			setupMock: func(m *mockApp.MockFindUsersPaginatedUseCase) {
				m.EXPECT().Execute(mock.MatchedBy(func(p domain.Pagination) bool {
					return p.LastID() != nil && *p.LastID() == 5 && p.Limit() == 25
				})).Return(paginatedUsers, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"users"`,
		},
		{
			testName:       "fail - invalid cursor string",
			urlQuery:       "?cursor=abc",
			setupMock:      func(m *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter cursor must be a positive integer",
		},
		{
			testName:       "fail - invalid cursor negative",
			urlQuery:       "?cursor=-5",
			setupMock:      func(m *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter cursor must be a positive integer",
		},
		{
			testName:       "fail - invalid limit string",
			urlQuery:       "?limit=abc",
			setupMock:      func(m *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter limit must be a positive integer",
		},
		{
			testName:       "fail - invalid limit negative",
			urlQuery:       "?limit=-10",
			setupMock:      func(m *mockApp.MockFindUsersPaginatedUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "parameter limit must be a positive integer",
		},
		{
			testName: "fail - usecase execution error",
			urlQuery: "",
			setupMock: func(m *mockApp.MockFindUsersPaginatedUseCase) {
				m.EXPECT().Execute(mock.Anything).Return(domain.PaginatedUsers{}, domain.ErrFindingUsers)
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
			rec := httptest.NewRecorder()

			controller.Handle(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedBody != "" {
				var responseMap map[string]any
				// If expectedBody is not a JSON object but a substring we expect in raw body, check that
				// otherwise unmarshal and check. Here let's just assert on the raw string for simplicity
				assert.Contains(t, rec.Body.String(), tc.expectedBody)
				_ = responseMap
			}
		})
	}
}

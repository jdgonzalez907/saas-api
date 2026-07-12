package controllers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"jdgonzalez907/users-api/internal/domain"
	"jdgonzalez907/users-api/internal/infrastructure/controllers"
	mockApp "jdgonzalez907/users-api/mocks/application"
)

func TestUserController_FindByID(t *testing.T) {
	now := time.Now()
	phone, _ := domain.NewPhone("+573001234567")
	email, _ := domain.NewEmail("john.doe@example.com")
	address, _ := domain.NewAddress("Calle 123", "Bogota", "Cundinamarca", "Colombia", nil, nil)
	birthDate, _ := domain.NewBirthDate(now.AddDate(-25, 0, 0))
	identification, _ := domain.NewIdentification(domain.IdType_CC, "12345678")

	userParams := domain.UserParams{
		ID:             1,
		Identification: identification,
		FirstName:      "John",
		LastName:       "Doe",
		Phone:          phone,
		Email:          &email,
		Address:        &address,
		BirthDate:      &birthDate,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	validUser, err := domain.NewUser(userParams)
	if err != nil {
		t.Fatalf("failed to create valid user: %v", err)
	}

	tests := []struct {
		name           string
		userIDParam    string
		setupMock      func(m *mockApp.MockFindUserByIdUseCase)
		expectedStatus int
		verifyResponse func(t *testing.T, body string)
	}{
		{
			name:        "success - user found",
			userIDParam: "1",
			setupMock: func(m *mockApp.MockFindUserByIdUseCase) {
				m.On("Execute", 1).Return(validUser, nil)
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body string) {
				var res domain.UserDTO
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Equal(t, validUser.ToDTO().ID, res.ID)
				assert.Equal(t, validUser.ToDTO().FirstName, res.FirstName)
				assert.Equal(t, validUser.ToDTO().LastName, res.LastName)
			},
		},
		{
			name:           "fail - non-numeric id format",
			userIDParam:    "abc",
			setupMock:      func(m *mockApp.MockFindUserByIdUseCase) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				var res controllers.ErrorResponse
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Contains(t, res.Message, "must be a positive integer")
			},
		},
		{
			name:           "fail - non-positive id",
			userIDParam:    "0",
			setupMock:      func(m *mockApp.MockFindUserByIdUseCase) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body string) {
				var res controllers.ErrorResponse
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Contains(t, res.Message, "must be a positive integer")
			},
		},
		{
			name:        "fail - internal server error",
			userIDParam: "1",
			setupMock: func(m *mockApp.MockFindUserByIdUseCase) {
				m.On("Execute", 1).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, body string) {
				var res controllers.ErrorResponse
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Equal(t, "internal server error", res.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := mockApp.NewMockFindUserByIdUseCase(t)
			tt.setupMock(mockUC)

			controller := controllers.NewUserController(mockUC)
			router := controllers.NewRouter(controller)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", tt.userIDParam), nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			tt.verifyResponse(t, w.Body.String())
		})
	}
}

func TestRespondWithJSON_NilData(t *testing.T) {
	w := httptest.NewRecorder()
	controllers.RespondWithJSON(w, http.StatusNoContent, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestParseRouteIntParam_MissingParam(t *testing.T) {
	// Request without a chi router context → URLParam returns ""
	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	_, err := controllers.ParseRouteIntParam(req, "id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is missing")
}

func TestRespondWithDomainError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		// Branch 1: known error, non-500 status → expose sentinel error message
		{"known error 404", domain.ErrUserNotFound, http.StatusNotFound, "user not found"},
		// Branch 2: known infra error, 500 status → hide internal details
		{"known infra error 500", domain.ErrCreatingUser, http.StatusInternalServerError, "internal server error"},
		// Branch 3: unknown error, not in map → default 500
		{"unknown error 500", errors.New("unexpected failure"), http.StatusInternalServerError, "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			controllers.RespondWithDomainError(w, tt.err)
			assert.Equal(t, tt.expectedStatus, w.Code)
			var res controllers.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, res.Message)
		})
	}
}

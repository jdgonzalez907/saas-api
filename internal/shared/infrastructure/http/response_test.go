package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	postsDomain "jdgonzalez907/saas-api/internal/posts/domain"
	sharedHttp "jdgonzalez907/saas-api/internal/shared/infrastructure/http"
	usersDomain "jdgonzalez907/saas-api/internal/users/domain"
)

func TestProtected(t *testing.T) {
	testCases := []struct {
		testName       string
		authHeader     string
		expectedStatus int
		expectedUserID int64
	}{
		{
			testName:       "success - plain numeric token",
			authHeader:     "42",
			expectedStatus: http.StatusOK,
			expectedUserID: 42,
		},
		{
			testName:       "success - bearer prefix",
			authHeader:     "Bearer 42",
			expectedStatus: http.StatusOK,
			expectedUserID: 42,
		},
		{
			testName:       "success - bearer prefix lowercase",
			authHeader:     "bearer 42",
			expectedStatus: http.StatusOK,
			expectedUserID: 42,
		},
		{
			testName:       "success - bearer prefix with spaces",
			authHeader:     "Bearer  42 ",
			expectedStatus: http.StatusOK,
			expectedUserID: 42,
		},
		{
			testName:       "fail - empty authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
		{
			testName:       "fail - non-numeric token",
			authHeader:     "abc",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
		{
			testName:       "fail - bearer with non-numeric token",
			authHeader:     "Bearer abc",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
		{
			testName:       "success - negative number is valid token",
			authHeader:     "-1",
			expectedStatus: http.StatusOK,
			expectedUserID: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var capturedUserID int64
			nextHandler := sharedHttp.Protected(func(w http.ResponseWriter, _ *http.Request, userID int64) {
				capturedUserID = userID
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()

			nextHandler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedUserID, capturedUserID)
			}
		})
	}
}

func TestResponseHelpers(t *testing.T) {
	t.Run("RespondWithJSON - nil data", func(t *testing.T) {
		rec := httptest.NewRecorder()
		sharedHttp.RespondWithJSON(rec, http.StatusNoContent, nil)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())
	})

	t.Run("RespondWithJSON - success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		data := map[string]string{"foo": "bar"}
		sharedHttp.RespondWithJSON(rec, http.StatusOK, data)
		assert.Equal(t, http.StatusOK, rec.Code)

		var res map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "bar", res["foo"])
	})

	t.Run("RespondWithDomainError - known 4xx error (users)", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		sharedHttp.RespondWithDomainError(rec, r, usersDomain.ErrUserNotFound)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var res map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, usersDomain.ErrUserNotFound.Error(), res["message"])
	})

	t.Run("RespondWithDomainError - known 4xx error (posts)", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		sharedHttp.RespondWithDomainError(rec, r, postsDomain.ErrPostNotFound)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var res map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, postsDomain.ErrPostNotFound.Error(), res["message"])
	})

	t.Run("RespondWithDomainError - known 500 error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		sharedHttp.RespondWithDomainError(rec, r, usersDomain.ErrCreatingUser)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var res map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "internal server error", res["message"])
	})

	t.Run("RespondWithDomainError - unknown error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		sharedHttp.RespondWithDomainError(rec, r, errors.New("unknown error happened"))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var res map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "internal server error", res["message"])
	})

	t.Run("ErrorLoggerMiddleware - logs internal server error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sharedHttp.RespondWithDomainError(w, r, usersDomain.ErrCreatingUser)
		})

		middleware := sharedHttp.ErrorLoggerMiddleware(nextHandler)
		middleware.ServeHTTP(rec, r)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ErrorLoggerMiddleware - does not log client error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sharedHttp.RespondWithDomainError(w, r, usersDomain.ErrUserNotFound)
		})

		middleware := sharedHttp.ErrorLoggerMiddleware(nextHandler)
		middleware.ServeHTTP(rec, r)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("JSONContentTypeMiddleware - sets header", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := sharedHttp.JSONContentTypeMiddleware(nextHandler)
		middleware.ServeHTTP(rec, r)

		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})
}

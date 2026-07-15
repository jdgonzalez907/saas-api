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

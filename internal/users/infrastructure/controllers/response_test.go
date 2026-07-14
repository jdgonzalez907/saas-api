package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"jdgonzalez907/saas-api/internal/users/domain"
	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"

	"github.com/stretchr/testify/assert"
)

func TestResponseHelpers(t *testing.T) {
	t.Run("RespondWithJSON - nil data", func(t *testing.T) {
		rec := httptest.NewRecorder()
		controllers.RespondWithJSON(rec, http.StatusNoContent, nil)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())
	})

	t.Run("RespondWithJSON - success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		data := map[string]string{"foo": "bar"}
		controllers.RespondWithJSON(rec, http.StatusOK, data)
		assert.Equal(t, http.StatusOK, rec.Code)

		var res map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "bar", res["foo"])
	})

	t.Run("RespondWithDomainError - known 4xx error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		controllers.RespondWithDomainError(rec, r, domain.ErrUserNotFound)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var res map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, domain.ErrUserNotFound.Error(), res["message"])
	})

	t.Run("RespondWithDomainError - known 500 error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		controllers.RespondWithDomainError(rec, r, domain.ErrCreatingUser)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var res map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "internal server error", res["message"])
	})

	t.Run("RespondWithDomainError - unknown error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		controllers.RespondWithDomainError(rec, r, errors.New("unknown error happened"))
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
			controllers.RespondWithDomainError(w, r, domain.ErrCreatingUser)
		})

		middleware := controllers.ErrorLoggerMiddleware(nextHandler)
		middleware.ServeHTTP(rec, r)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ErrorLoggerMiddleware - does not log client error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.RespondWithDomainError(w, r, domain.ErrUserNotFound)
		})

		middleware := controllers.ErrorLoggerMiddleware(nextHandler)
		middleware.ServeHTTP(rec, r)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

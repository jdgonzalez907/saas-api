package http_test

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
)

func TestParseRouteInt64Param(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/42", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "42")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		val, err := sharedHttp.ParseRouteInt64Param(req, "id")
		assert.NoError(t, err)
		assert.Equal(t, int64(42), val)
	})

	t.Run("fail - missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", nil)
		val, err := sharedHttp.ParseRouteInt64Param(req, "id")
		assert.ErrorIs(t, err, sharedHttp.ErrInvalidRouteParam)
		assert.Equal(t, int64(0), val)
	})

	t.Run("fail - invalid negative/zero", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/-42", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "-42")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		val, err := sharedHttp.ParseRouteInt64Param(req, "id")
		assert.ErrorIs(t, err, sharedHttp.ErrInvalidRouteParam)
		assert.Equal(t, int64(0), val)
	})
}

func TestParseQueryInt64Param(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users?last_id=42", nil)
		val, err := sharedHttp.ParseQueryInt64Param(req, "last_id")
		assert.NoError(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, int64(42), *val)
	})

	t.Run("success - empty", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", nil)
		val, err := sharedHttp.ParseQueryInt64Param(req, "last_id")
		assert.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("fail - invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users?last_id=abc", nil)
		val, err := sharedHttp.ParseQueryInt64Param(req, "last_id")
		assert.ErrorIs(t, err, sharedHttp.ErrInvalidQueryParam)
		assert.Nil(t, val)
	})
}

func TestParseQueryInt32Param(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users?limit=42", nil)
		val, err := sharedHttp.ParseQueryInt32Param(req, "limit")
		assert.NoError(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, int32(42), *val)
	})

	t.Run("success - empty", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users", nil)
		val, err := sharedHttp.ParseQueryInt32Param(req, "limit")
		assert.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("fail - invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users?limit=abc", nil)
		val, err := sharedHttp.ParseQueryInt32Param(req, "limit")
		assert.ErrorIs(t, err, sharedHttp.ErrInvalidQueryParam)
		assert.Nil(t, val)
	})
}

func TestDecodeJSON(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
	}

	t.Run("success", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name":"test"}`)
		req := httptest.NewRequest("POST", "/", body)
		var dst testStruct
		err := sharedHttp.DecodeJSON(req, &dst)
		assert.NoError(t, err)
		assert.Equal(t, "test", dst.Name)
	})

	t.Run("fail - nil body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", nil)
		req.Body = nil
		var dst testStruct
		err := sharedHttp.DecodeJSON(req, &dst)
		assert.ErrorIs(t, err, sharedHttp.ErrInvalidRequestBody)
	})

	t.Run("fail - invalid syntax", func(t *testing.T) {
		body := bytes.NewBufferString(`{invalid}`)
		req := httptest.NewRequest("POST", "/", body)
		var dst testStruct
		err := sharedHttp.DecodeJSON(req, &dst)
		assert.ErrorIs(t, err, sharedHttp.ErrInvalidRequestBody)
	})
}

func TestParseQueryTimeParam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/posts?last_published_at=2026-07-14T19:00:00Z", nil)
		val, err := sharedHttp.ParseQueryTimeParam(req, "last_published_at")
		assert.NoError(t, err)
		assert.NotNil(t, val)
		expected, _ := time.Parse(time.RFC3339, "2026-07-14T19:00:00Z")
		assert.True(t, expected.Equal(*val))
	})

	t.Run("success - empty", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/posts", nil)
		val, err := sharedHttp.ParseQueryTimeParam(req, "last_published_at")
		assert.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("fail - invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/posts?last_published_at=abc", nil)
		val, err := sharedHttp.ParseQueryTimeParam(req, "last_published_at")
		assert.ErrorIs(t, err, sharedHttp.ErrInvalidQueryParam)
		assert.Nil(t, val)
	})
}

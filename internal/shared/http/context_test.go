package http_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	sharedHttp "jdgonzalez907/saas-api/internal/shared/http"
)

func TestGetAuthenticatedUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", int64(42))
		id, err := sharedHttp.GetAuthenticatedUserID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("fail - not in context", func(t *testing.T) {
		id, err := sharedHttp.GetAuthenticatedUserID(context.Background())
		assert.ErrorIs(t, err, sharedHttp.ErrUnauthenticated)
		assert.Equal(t, int64(0), id)
	})

	t.Run("fail - wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "userID", "42")
		id, err := sharedHttp.GetAuthenticatedUserID(ctx)
		assert.ErrorIs(t, err, sharedHttp.ErrUnauthenticated)
		assert.Equal(t, int64(0), id)
	})
}

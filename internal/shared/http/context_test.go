package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthenticatedUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := WithUserID(context.Background(), int64(42))
		id, err := GetAuthenticatedUserID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("fail - not in context", func(t *testing.T) {
		id, err := GetAuthenticatedUserID(context.Background())
		assert.ErrorIs(t, err, ErrUnauthenticated)
		assert.Equal(t, int64(0), id)
	})

	t.Run("fail - wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "42")
		id, err := GetAuthenticatedUserID(ctx)
		assert.ErrorIs(t, err, ErrUnauthenticated)
		assert.Equal(t, int64(0), id)
	})
}

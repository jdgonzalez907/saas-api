package http

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "userID"

var ErrUnauthenticated = errors.New("unauthenticated user")

func GetAuthenticatedUserID(ctx context.Context) (int64, error) {
	val, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0, ErrUnauthenticated
	}
	return val, nil
}

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

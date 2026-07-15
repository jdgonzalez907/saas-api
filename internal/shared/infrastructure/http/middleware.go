package http

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func ErrorLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		holder := &ErrorHolder{}
		ctx := context.WithValue(r.Context(), ErrorHolderKey, holder)

		next.ServeHTTP(w, r.WithContext(ctx))

		if holder.Err != nil {
			status, _ := statusFromDomainError(holder.Err)
			if status == http.StatusInternalServerError {
				reqID := middleware.GetReqID(r.Context())
				log.Printf("[ERROR] RequestID: %s - Internal Server Error: %v", reqID, holder.Err)
			}
		}
	})
}

type ProtectedHandlerFunc func(w http.ResponseWriter, r *http.Request, userID int64)

func Protected(next ProtectedHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithDomainError(w, r, ErrUnauthenticated)
			return
		}

		token := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = authHeader[7:]
		}

		userID, err := strconv.ParseInt(strings.TrimSpace(token), 10, 64)
		if err != nil {
			RespondWithDomainError(w, r, ErrUnauthenticated)
			return
		}

		next(w, r, userID)
	})
}

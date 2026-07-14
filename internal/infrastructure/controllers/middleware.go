package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type contextKey string

const errorHolderKey contextKey = "error_holder"

type ErrorHolder struct {
	Err error
}

func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func ErrorLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		holder := &ErrorHolder{}
		ctx := context.WithValue(r.Context(), errorHolderKey, holder)

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

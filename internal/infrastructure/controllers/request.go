package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var (
	ErrInvalidRouteParam   = errors.New("invalid route parameter")
	ErrInvalidRequestBody  = errors.New("invalid request body")
)

func ParseRouteIntParam(r *http.Request, paramName string) (int, error) {
	valStr := chi.URLParam(r, paramName)
	if valStr == "" {
		return 0, fmt.Errorf("%w: parameter %s is missing", ErrInvalidRouteParam, paramName)
	}

	val, err := strconv.Atoi(valStr)
	if err != nil || val <= 0 {
		return 0, fmt.Errorf("%w: parameter %s must be a positive integer", ErrInvalidRouteParam, paramName)
	}

	return val, nil
}

func ParseJSONBody(r *http.Request, dst any) error {
	if r.Body == nil {
		return fmt.Errorf("%w: request body is required", ErrInvalidRequestBody)
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidRequestBody, err.Error())
	}
	return nil
}

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
	ErrInvalidRouteParam = errors.New("invalid route parameter")
	ErrInvalidQueryParam = errors.New("invalid query parameter")
	ErrInvalidRequestBody = errors.New("invalid request body")
)

func ParseRouteInt64Param(r *http.Request, paramName string) (int64, error) {
	valStr := chi.URLParam(r, paramName)
	if valStr == "" {
		return 0, fmt.Errorf("%w: parameter %s is missing", ErrInvalidRouteParam, paramName)
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil || val <= 0 {
		return 0, fmt.Errorf("%w: parameter %s must be a positive integer", ErrInvalidRouteParam, paramName)
	}

	return val, nil
}

func ParseQueryInt64Param(r *http.Request, paramName string) (*int64, error) {
	valStr := r.URL.Query().Get(paramName)
	if valStr == "" {
		return nil, nil
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil || val <= 0 {
		return nil, fmt.Errorf("%w: parameter %s must be a positive integer", ErrInvalidQueryParam, paramName)
	}

	return &val, nil
}

func ParseQueryInt32Param(r *http.Request, paramName string) (*int32, error) {
	valStr := r.URL.Query().Get(paramName)
	if valStr == "" {
		return nil, nil
	}

	val, err := strconv.ParseInt(valStr, 10, 32)
	if err != nil || val <= 0 {
		return nil, fmt.Errorf("%w: parameter %s must be a positive integer", ErrInvalidQueryParam, paramName)
	}

	val32 := int32(val)
	return &val32, nil
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

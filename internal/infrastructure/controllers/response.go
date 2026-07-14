package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"jdgonzalez907/users-api/internal/domain"
)

const InternalServerErrorMessage = "internal server error"

type ErrorResponse struct {
	Message string `json:"message"`
}

var domainErrorStatus = map[error]int{
	domain.ErrUserNotFound: http.StatusNotFound,

	domain.ErrUserIDAlreadyExists:    http.StatusConflict,
	domain.ErrUserPhoneAlreadyExists: http.StatusConflict,
	domain.ErrUserEmailAlreadyExists: http.StatusConflict,

	domain.ErrInvalidUserID:             http.StatusBadRequest,
	domain.ErrInvalidFirstName:          http.StatusBadRequest,
	domain.ErrInvalidLastName:           http.StatusBadRequest,
	domain.ErrInvalidPhone:              http.StatusBadRequest,
	domain.ErrInvalidEmail:              http.StatusBadRequest,
	domain.ErrUserUnderage:              http.StatusBadRequest,
	domain.ErrInvalidBirthDateFormat:    http.StatusBadRequest,
	domain.ErrInvalidIdentificationType:   http.StatusBadRequest,
	domain.ErrInvalidIdentificationNumber: http.StatusBadRequest,
	domain.ErrInvalidStreet:             http.StatusBadRequest,
	domain.ErrInvalidCity:               http.StatusBadRequest,
	domain.ErrInvalidState:              http.StatusBadRequest,
	domain.ErrInvalidCountry:            http.StatusBadRequest,
	domain.ErrInvalidPostalCode:         http.StatusBadRequest,
	domain.ErrInvalidDescription:        http.StatusBadRequest,
	domain.ErrInvalidPaginationLimit:    http.StatusBadRequest,
	domain.ErrInvalidPaginationCursor:   http.StatusBadRequest,

	domain.ErrCreatingUser:                    http.StatusInternalServerError,
	domain.ErrUpdatingUserPersonalInformation: http.StatusInternalServerError,
	domain.ErrUpdatingUserPhone:               http.StatusInternalServerError,
	domain.ErrUpdatingUserEmail:               http.StatusInternalServerError,
	domain.ErrDeletingUser:                    http.StatusInternalServerError,
	domain.ErrFindingUsers:                    http.StatusInternalServerError,
	domain.ErrFindingUserByID:                 http.StatusInternalServerError,
}

func statusFromDomainError(err error) (int, string) {
	for domainErr, status := range domainErrorStatus {
		if errors.Is(err, domainErr) {
			msg := domainErr.Error()
			if status == http.StatusInternalServerError {
				msg = InternalServerErrorMessage
			}
			return status, msg
		}
	}
	return http.StatusInternalServerError, InternalServerErrorMessage
}

func RespondWithJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	if data == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, ErrorResponse{Message: message})
}

func RespondWithDomainError(w http.ResponseWriter, r *http.Request, err error) {
	status, msg := statusFromDomainError(err)
	if holder, ok := r.Context().Value(errorHolderKey).(*ErrorHolder); ok {
		holder.Err = err
	}
	RespondWithError(w, status, msg)
}

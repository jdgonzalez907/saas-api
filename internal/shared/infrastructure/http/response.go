package http

import (
	"encoding/json"
	"errors"
	"net/http"

	postsDomain "jdgonzalez907/saas-api/internal/posts/domain"
	usersDomain "jdgonzalez907/saas-api/internal/users/domain"
)

const InternalServerErrorMessage = "internal server error"

type ErrorResponse struct {
	Message string `json:"message"`
}

type ErrorHolder struct {
	Err error
}

type contextKey string

const ErrorHolderKey contextKey = "error_holder"

var domainErrorStatus = map[error]int{
	ErrUnauthenticated: http.StatusUnauthorized,

	usersDomain.ErrUserNotFound: http.StatusNotFound,

	usersDomain.ErrUserOwnershipMismatch: http.StatusForbidden,

	usersDomain.ErrUserIDAlreadyExists:    http.StatusConflict,
	usersDomain.ErrUserPhoneAlreadyExists: http.StatusConflict,
	usersDomain.ErrUserEmailAlreadyExists: http.StatusConflict,

	usersDomain.ErrInvalidUserID:               http.StatusBadRequest,
	usersDomain.ErrInvalidFirstName:            http.StatusBadRequest,
	usersDomain.ErrInvalidLastName:             http.StatusBadRequest,
	usersDomain.ErrInvalidPhone:                http.StatusBadRequest,
	usersDomain.ErrInvalidEmail:                http.StatusBadRequest,
	usersDomain.ErrUserUnderage:                http.StatusBadRequest,
	usersDomain.ErrInvalidBirthDateFormat:      http.StatusBadRequest,
	usersDomain.ErrInvalidIdentificationType:   http.StatusBadRequest,
	usersDomain.ErrInvalidIdentificationNumber: http.StatusBadRequest,
	usersDomain.ErrInvalidStreet:               http.StatusBadRequest,
	usersDomain.ErrInvalidCity:                 http.StatusBadRequest,
	usersDomain.ErrInvalidState:                http.StatusBadRequest,
	usersDomain.ErrInvalidCountry:              http.StatusBadRequest,
	usersDomain.ErrInvalidPostalCode:           http.StatusBadRequest,
	usersDomain.ErrInvalidDescription:          http.StatusBadRequest,
	usersDomain.ErrInvalidPaginationLimit:      http.StatusBadRequest,
	usersDomain.ErrInvalidPaginationCursor:     http.StatusBadRequest,

	usersDomain.ErrCreatingUser:                http.StatusInternalServerError,
	usersDomain.ErrChangingPersonalInformation: http.StatusInternalServerError,
	usersDomain.ErrChangingPhone:               http.StatusInternalServerError,
	usersDomain.ErrChangingEmail:               http.StatusInternalServerError,
	usersDomain.ErrDeletingUser:                http.StatusInternalServerError,
	usersDomain.ErrFindingUsers:                http.StatusInternalServerError,
	usersDomain.ErrFindingUserByID:             http.StatusInternalServerError,

	postsDomain.ErrPostNotFound: http.StatusNotFound,

	postsDomain.ErrPostOwnershipMismatch: http.StatusForbidden,

	postsDomain.ErrPostIDAlreadyExists: http.StatusConflict,

	postsDomain.ErrInvalidPostID:                    http.StatusBadRequest,
	postsDomain.ErrInvalidPostStatus:                http.StatusBadRequest,
	postsDomain.ErrInvalidAuthorID:                  http.StatusBadRequest,
	postsDomain.ErrDraftCannotHavePublicationDate:   http.StatusBadRequest,
	postsDomain.ErrPublishedMustHavePublicationDate: http.StatusBadRequest,
	postsDomain.ErrEmptyPostTitle:                   http.StatusBadRequest,
	postsDomain.ErrOrphanBlock:                      http.StatusBadRequest,
	postsDomain.ErrInvalidPaginationLimit:           http.StatusBadRequest,
	postsDomain.ErrInvalidPaginationCursor:          http.StatusBadRequest,

	postsDomain.ErrCreatingPost: http.StatusInternalServerError,
	postsDomain.ErrFindingPost:  http.StatusInternalServerError,
	postsDomain.ErrChangingPost: http.StatusInternalServerError,
	postsDomain.ErrDeletingPost: http.StatusInternalServerError,
	postsDomain.ErrFindingPosts: http.StatusInternalServerError,
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
	if holder, ok := r.Context().Value(ErrorHolderKey).(*ErrorHolder); ok {
		holder.Err = err
	}
	RespondWithError(w, status, msg)
}

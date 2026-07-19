package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jdgonzalez907/saas-api/internal/users/domain"
	mock_domain "github.com/jdgonzalez907/saas-api/mocks/users/domain"
	"github.com/stretchr/testify/mock"
)

func TestChangeEmail_Execute(t *testing.T) {
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	email := validEmail(t)
	newEmail, _ := domain.NewEmail("new@example.com")
	now := time.Now()

	tests := []struct {
		name       string
		setup      func(t *testing.T, repo *mock_domain.MockUserRepository)
		executedBy int64
		email      domain.Email
		wantUser   bool
		wantErr    error
	}{
		{
			name: "success - change to new email",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByEmail", mock.Anything, newEmail).Return(nil, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			executedBy: 1,
			email:      newEmail,
			wantUser:   true,
			wantErr:    nil,
		},
		{
			name: "success - same email",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByEmail", mock.Anything, *email).Return(user, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			executedBy: 1,
			email:      *email,
			wantUser:   true,
			wantErr:    nil,
		},
		{
			name: "error - user not found",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			executedBy: 1,
			email:      newEmail,
			wantUser:   false,
			wantErr:    domain.ErrChangeEmail,
		},
		{
			name: "error - email already exists",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				other := mustNewUserWithID(t, 99, &newEmail, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByEmail", mock.Anything, newEmail).Return(other, nil)
			},
			executedBy: 1,
			email:      newEmail,
			wantUser:   false,
			wantErr:    domain.ErrChangeEmail,
		},
		{
			name: "error - FindByID fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))
			},
			executedBy: 1,
			email:      newEmail,
			wantUser:   false,
			wantErr:    domain.ErrChangeEmail,
		},
		{
			name: "error - FindByEmail fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByEmail", mock.Anything, newEmail).Return(nil, errors.New("database error"))
			},
			executedBy: 1,
			email:      newEmail,
			wantUser:   false,
			wantErr:    domain.ErrChangeEmail,
		},
		{
			name: "error - Update fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByEmail", mock.Anything, newEmail).Return(nil, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			executedBy: 1,
			email:      newEmail,
			wantUser:   false,
			wantErr:    domain.ErrChangeEmail,
		},
		{
			name: "error - ChangeEmail fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(99)).Return(user, nil)
				repo.On("FindByEmail", mock.Anything, newEmail).Return(nil, nil)
			},
			executedBy: 99,
			email:      newEmail,
			wantUser:   false,
			wantErr:    domain.ErrChangeEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockUserRepository(t)
			tt.setup(t, repo)
			uc := NewChangeEmail(repo)
			got, err := uc.Execute(context.Background(), tt.executedBy, tt.email)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantUser && got == nil {
				t.Errorf("Execute() returned nil user, want non-nil")
			}

			if !tt.wantUser && got != nil {
				t.Errorf("Execute() returned non-nil user, want nil")
			}
		})
	}
}

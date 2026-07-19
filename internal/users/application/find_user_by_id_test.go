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

func TestFindUserByID_Execute(t *testing.T) {
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	email := validEmail(t)
	now := time.Now()

	tests := []struct {
		name      string
		setup     func(t *testing.T, repo *mock_domain.MockUserRepository)
		id        int64
		wantUser  bool
		wantErr   error
	}{
		{
			name: "success",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
			},
			id:       1,
			wantUser: true,
			wantErr:  nil,
		},
		{
			name: "error - user not found",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(99)).Return(nil, domain.ErrUserNotFound)
			},
			id:       99,
			wantUser: false,
			wantErr:  domain.ErrFindUserByID,
		},
		{
			name: "error - repository error",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, errors.New("database connection failed"))
			},
			id:       1,
			wantUser: false,
			wantErr:  domain.ErrFindUserByID,
		},
		{
			name: "error - user nil",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(2)).Return(nil, nil)
			},
			id:       2,
			wantUser: false,
			wantErr:  domain.ErrFindUserByID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockUserRepository(t)
			tt.setup(t, repo)
			uc := NewFindUserByID(repo)
			got, err := uc.Execute(context.Background(), tt.id)

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
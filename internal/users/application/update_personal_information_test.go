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

func TestUpdatePersonalInformation_Execute(t *testing.T) {
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	email := validEmail(t)
	newPI, _ := domain.NewPersonalInformation(domain.DNITypeCE, "9876543210", "Jane", "DoeSmith", time.Date(1985, 6, 15, 0, 0, 0, 0, time.UTC))
	now := time.Now()

	tests := []struct {
		name                string
		setup               func(t *testing.T, repo *mock_domain.MockUserRepository)
		executedBy          int64
		personalInformation domain.PersonalInformation
		wantUser            bool
		wantErr             error
	}{
		{
			name: "success - update personal information",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			executedBy:          1,
			personalInformation: newPI,
			wantUser:            true,
			wantErr:             nil,
		},
		{
			name: "success - same personal information",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			executedBy:          1,
			personalInformation: pi,
			wantUser:            true,
			wantErr:             nil,
		},
		{
			name: "error - user not found",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			executedBy:          1,
			personalInformation: newPI,
			wantUser:            false,
			wantErr:             domain.ErrUpdatePersonalInformation,
		},
		{
			name: "error - unauthorized modification",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(2)).Return(user, nil)
			},
			executedBy:          2,
			personalInformation: newPI,
			wantUser:            false,
			wantErr:             domain.ErrUpdatePersonalInformation,
		},
		{
			name: "error - FindByID fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))
			},
			executedBy:          1,
			personalInformation: newPI,
			wantUser:            false,
			wantErr:             domain.ErrUpdatePersonalInformation,
		},
		{
			name: "error - Update fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, email, pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			executedBy:          1,
			personalInformation: newPI,
			wantUser:            false,
			wantErr:             domain.ErrUpdatePersonalInformation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockUserRepository(t)
			tt.setup(t, repo)
			uc := NewUpdatePersonalInformation(repo)
			got, err := uc.Execute(context.Background(), tt.executedBy, tt.personalInformation)

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

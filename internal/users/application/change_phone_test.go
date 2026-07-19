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

func mustNewUserWithID(t *testing.T, id int64, email *domain.Email, pi domain.PersonalInformation, phone domain.Phone, createdAt, updatedAt time.Time) *domain.User {
	t.Helper()
	u, err := domain.NewWithID(id, email, pi, phone, createdAt, updatedAt)
	if err != nil {
		t.Fatalf("mustNewUserWithID() error = %v", err)
	}
	return u
}

func TestChangePhone_Execute(t *testing.T) {
	pi := validPersonalInformation(t)
	phone := validPhone(t)
	newPhone, _ := domain.NewPhone("1", "2025551234")
	now := time.Now()

	tests := []struct {
		name       string
		setup      func(t *testing.T, repo *mock_domain.MockUserRepository)
		executedBy int64
		phone      domain.Phone
		wantUser   bool
		wantErr    error
	}{
		{
			name: "success - change to new phone",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByPhone", mock.Anything, newPhone).Return(nil, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			executedBy: 1,
			phone:      newPhone,
			wantUser:   true,
			wantErr:    nil,
		},
		{
			name: "success - same phone",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByPhone", mock.Anything, phone).Return(user, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			executedBy: 1,
			phone:      phone,
			wantUser:   true,
			wantErr:    nil,
		},
		{
			name: "error - user not found",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			executedBy: 1,
			phone:      newPhone,
			wantUser:   false,
			wantErr:    domain.ErrChangePhone,
		},
		{
			name: "error - phone already exists",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now)
				other := mustNewUserWithID(t, 99, validEmail(t), pi, newPhone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByPhone", mock.Anything, newPhone).Return(other, nil)
			},
			executedBy: 1,
			phone:      newPhone,
			wantUser:   false,
			wantErr:    domain.ErrChangePhone,
		},
		{
			name: "error - FindByID fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByID", mock.Anything, int64(1)).Return(nil, errors.New("database error"))
			},
			executedBy: 1,
			phone:      newPhone,
			wantUser:   false,
			wantErr:    domain.ErrChangePhone,
		},
		{
			name: "error - FindByPhone fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByPhone", mock.Anything, newPhone).Return(nil, errors.New("database error"))
			},
			executedBy: 1,
			phone:      newPhone,
			wantUser:   false,
			wantErr:    domain.ErrChangePhone,
		},
		{
			name: "error - Update fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(1)).Return(user, nil)
				repo.On("FindByPhone", mock.Anything, newPhone).Return(nil, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			executedBy: 1,
			phone:      newPhone,
			wantUser:   false,
			wantErr:    domain.ErrChangePhone,
		},
		{
			name: "error - ChangePhone fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				user := mustNewUserWithID(t, 1, validEmail(t), pi, phone, now, now)
				repo.On("FindByID", mock.Anything, int64(99)).Return(user, nil)
				repo.On("FindByPhone", mock.Anything, newPhone).Return(nil, nil)
			},
			executedBy: 99,
			phone:      newPhone,
			wantUser:   false,
			wantErr:    domain.ErrChangePhone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockUserRepository(t)
			tt.setup(t, repo)
			uc := NewChangePhone(repo)
			got, err := uc.Execute(context.Background(), tt.executedBy, tt.phone)

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

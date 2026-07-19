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

func mustNow() time.Time {
	return time.Now()
}

func validPersonalInformation(t *testing.T) domain.PersonalInformation {
	t.Helper()
	pi, err := domain.NewPersonalInformation(domain.DNITypeCC, "1234567890", "John", "DoeSmith", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("validPersonalInformation() error = %v", err)
	}
	return pi
}

func validPhone(t *testing.T) domain.Phone {
	t.Helper()
	p, err := domain.NewPhone("57", "3001234567")
	if err != nil {
		t.Fatalf("validPhone() error = %v", err)
	}
	return p
}

func validEmail(t *testing.T) *domain.Email {
	t.Helper()
	e, err := domain.NewEmail("john@example.com")
	if err != nil {
		t.Fatalf("validEmail() error = %v", err)
	}
	return &e
}

func mustNewUser(t *testing.T, email *domain.Email, pi domain.PersonalInformation, phone domain.Phone) *domain.User {
	t.Helper()
	u, err := domain.New(email, pi, phone)
	if err != nil {
		t.Fatalf("mustNewUser() error = %v", err)
	}
	return u
}

func TestCreateUser_Execute(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, repo *mock_domain.MockUserRepository)
		user     *domain.User
		wantUser bool
		wantErr  error
	}{
		{
			name: "success - without email",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByPhone", mock.Anything, mock.Anything).Return(nil, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			user:     mustNewUser(t, nil, validPersonalInformation(t), validPhone(t)),
			wantUser: true,
			wantErr:  nil,
		},
		{
			name: "success - with email",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByPhone", mock.Anything, mock.Anything).Return(nil, nil)
				repo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			user:     mustNewUser(t, validEmail(t), validPersonalInformation(t), validPhone(t)),
			wantUser: true,
			wantErr:  nil,
		},
		{
			name: "error - phone already exists",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				existing := mustNewUserWithID(t, 1, validEmail(t), validPersonalInformation(t), validPhone(t), mustNow(), mustNow())
				repo.On("FindByPhone", mock.Anything, mock.Anything).Return(existing, nil)
			},
			user:    mustNewUser(t, nil, validPersonalInformation(t), validPhone(t)),
			wantErr: domain.ErrCreateUser,
		},
		{
			name: "error - email already exists",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				existing := mustNewUserWithID(t, 1, validEmail(t), validPersonalInformation(t), validPhone(t), mustNow(), mustNow())
				repo.On("FindByPhone", mock.Anything, mock.Anything).Return(nil, nil)
				repo.On("FindByEmail", mock.Anything, mock.Anything).Return(existing, nil)
			},
			user:    mustNewUser(t, validEmail(t), validPersonalInformation(t), validPhone(t)),
			wantErr: domain.ErrCreateUser,
		},
		{
			name: "error - FindByPhone fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByPhone", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			user:    mustNewUser(t, nil, validPersonalInformation(t), validPhone(t)),
			wantErr: domain.ErrCreateUser,
		},
		{
			name: "error - FindByEmail fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByPhone", mock.Anything, mock.Anything).Return(nil, nil)
				repo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			user:    mustNewUser(t, validEmail(t), validPersonalInformation(t), validPhone(t)),
			wantErr: domain.ErrCreateUser,
		},
		{
			name: "error - Create fails",
			setup: func(t *testing.T, repo *mock_domain.MockUserRepository) {
				repo.On("FindByPhone", mock.Anything, mock.Anything).Return(nil, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			user:    mustNewUser(t, nil, validPersonalInformation(t), validPhone(t)),
			wantErr: domain.ErrCreateUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock_domain.NewMockUserRepository(t)
			tt.setup(t, repo)
			uc := NewCreateUser(repo)
			got, err := uc.Execute(context.Background(), tt.user)

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



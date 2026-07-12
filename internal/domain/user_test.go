package domain_test

import (
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewUser(t *testing.T) {
	identification := domain.Identification{Type: domain.IdType_CC, Number: "123456789"}
	email := domain.Email{Value: "name@domain.com"}
	phone := domain.Phone{Value: "123456789"}
	address := domain.Address{
		Street:      "123 Main St",
		PostalCode:  nil,
		City:        "New York",
		State:       "NY",
		Country:     "USA",
		Description: nil,
	}
	birthDate := domain.BirthDate{Value: time.Now()}

	testCases := []struct {
		testName       string
		input          domain.User
		expectedError  error
		expectedOutput *domain.User
	}{
		{
			testName: "create user",
			input: domain.User{
				ID:             "1",
				Identification: identification,
				FirstName:      "John",
				LastName:       "Doe",
				Phone:          phone,
				Email:          &email,
				Address:        &address,
				BirthDate:      &birthDate,
			},
			expectedError: nil,
			expectedOutput: &domain.User{
				ID:             "1",
				Identification: identification,
				FirstName:      "John",
				LastName:       "Doe",
				Phone:          phone,
				Email:          &email,
				Address:        &address,
				BirthDate:      &birthDate,
			},
		},
		{
			testName: "fail to create user with empty id",
			input: domain.User{
				ID:             "",
				Identification: identification,
				FirstName:      "John",
				LastName:       "Doe",
				Phone:          phone,
				Email:          &email,
				Address:        &address,
				BirthDate:      &birthDate,
			},
			expectedError:  domain.ErrInvalidUserID,
			expectedOutput: nil,
		},
		{
			testName: "fail to create user with empty first name",
			input: domain.User{
				ID:             "1",
				Identification: identification,
				FirstName:      "",
				LastName:       "Doe",
				Phone:          phone,
				Email:          &email,
				Address:        &address,
				BirthDate:      &birthDate,
			},
			expectedError:  domain.ErrInvalidFirstName,
			expectedOutput: nil,
		},
		{
			testName: "fail to create user with empty last name",
			input: domain.User{
				ID:             "1",
				Identification: identification,
				FirstName:      "John",
				LastName:       "",
				Phone:          phone,
				Email:          &email,
				Address:        &address,
				BirthDate:      &birthDate,
			},
			expectedError:  domain.ErrInvalidLastName,
			expectedOutput: nil,
		},
	}

	for _, testCase := range testCases {
		user, err := domain.NewUser(
			testCase.input.ID,
			testCase.input.Identification,
			testCase.input.FirstName,
			testCase.input.LastName,
			testCase.input.Phone,
			testCase.input.Email,
			testCase.input.Address,
			testCase.input.BirthDate,
		)
		if err != testCase.expectedError {
			t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
		}
		if !isEqual(user, testCase.expectedOutput) {
			t.Errorf("expected user: %v, got %v", testCase.expectedOutput, user)
		}
	}
}

func isEqual(user1, user2 *domain.User) bool {
	if user1 == nil && user2 == nil {
		return true
	}
	if user1 == nil || user2 == nil {
		return false
	}
	return user1.ID == user2.ID &&
		user1.Identification == user2.Identification &&
		user1.FirstName == user2.FirstName &&
		user1.LastName == user2.LastName &&
		user1.Phone == user2.Phone &&
		user1.Email == user2.Email &&
		user1.Address == user2.Address &&
		user1.BirthDate == user2.BirthDate
}

func TestNewUserWithoutId(t *testing.T) {
	identification := domain.Identification{Type: domain.IdType_CC, Number: "123456789"}
	email := domain.Email{Value: "name@domain.com"}
	phone := domain.Phone{Value: "123456789"}
	address := domain.Address{
		Street:      "123 Main St",
		PostalCode:  nil,
		City:        "New York",
		State:       "NY",
		Country:     "USA",
		Description: nil,
	}
	birthDate := domain.BirthDate{Value: time.Now()}

	user, err := domain.NewUserWithoutId(
		identification,
		"John",
		"Doe",
		phone,
		&email,
		&address,
		&birthDate,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be not nil")
	}

	if user.ID == "" {
		t.Error("expected generated ID to be not empty")
	}

	_, err = domain.NewUserWithoutId(
		identification,
		"",
		"Doe",
		phone,
		&email,
		&address,
		&birthDate,
	)
	if err != domain.ErrInvalidFirstName {
		t.Errorf("expected ErrInvalidFirstName, got %v", err)
	}
}

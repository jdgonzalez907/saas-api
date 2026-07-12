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
	now := time.Now()

	testCases := []struct {
		testName       string
		id             int
		firstName      string
		lastName       string
		expectedError  error
		expectedOutput func(int) *domain.User
	}{
		{
			testName:      "create user",
			id:            1,
			firstName:     "John",
			lastName:      "Doe",
			expectedError: nil,
			expectedOutput: func(id int) *domain.User {
				u, _ := domain.NewUser(id, identification, "John", "Doe", phone, &email, &address, &birthDate, now, now)
				return u
			},
		},
		{
			testName:      "fail to create user with empty id (less than 0)",
			id:            -1,
			firstName:     "John",
			lastName:      "Doe",
			expectedError: domain.ErrInvalidUserID,
			expectedOutput: func(id int) *domain.User {
				return nil
			},
		},
		{
			testName:      "fail to create user with empty first name",
			id:            1,
			firstName:     "",
			lastName:      "Doe",
			expectedError: domain.ErrInvalidFirstName,
			expectedOutput: func(id int) *domain.User {
				return nil
			},
		},
		{
			testName:      "fail to create user with empty last name",
			id:            1,
			firstName:     "John",
			lastName:      "",
			expectedError: domain.ErrInvalidLastName,
			expectedOutput: func(id int) *domain.User {
				return nil
			},
		},
	}

	for _, testCase := range testCases {
		user, err := domain.NewUser(
			testCase.id,
			identification,
			testCase.firstName,
			testCase.lastName,
			phone,
			&email,
			&address,
			&birthDate,
			now,
			now,
		)
		if err != testCase.expectedError {
			t.Errorf("%s: expected error: %v, got %v", testCase.testName, testCase.expectedError, err)
		}
		expected := testCase.expectedOutput(testCase.id)
		if !isEqual(user, expected) {
			t.Errorf("%s: expected user: %v, got %v", testCase.testName, expected, user)
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
		user1.BirthDate == user2.BirthDate &&
		user1.CreatedAt.Equal(user2.CreatedAt) &&
		user1.UpdatedAt.Equal(user2.UpdatedAt)
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

	if user.ID != 0 {
		t.Errorf("expected generated ID to be 0, got %d", user.ID)
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

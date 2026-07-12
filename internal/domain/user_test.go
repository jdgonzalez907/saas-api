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
		testName      string
		id            int
		firstName     string
		lastName      string
		expectedError error
	}{
		{
			testName:      "create user",
			id:            1,
			firstName:     "John",
			lastName:      "Doe",
			expectedError: nil,
		},
		{
			testName:      "fail to create user with empty id (less than 0)",
			id:            -1,
			firstName:     "John",
			lastName:      "Doe",
			expectedError: domain.ErrInvalidUserID,
		},
		{
			testName:      "fail to create user with empty first name",
			id:            1,
			firstName:     "",
			lastName:      "Doe",
			expectedError: domain.ErrInvalidFirstName,
		},
		{
			testName:      "fail to create user with empty last name",
			id:            1,
			firstName:     "John",
			lastName:      "",
			expectedError: domain.ErrInvalidLastName,
		},
	}

	for _, testCase := range testCases {
		params := domain.UserParams{
			ID:             testCase.id,
			Identification: identification,
			FirstName:      testCase.firstName,
			LastName:       testCase.lastName,
			Phone:          phone,
			Email:          &email,
			Address:        &address,
			BirthDate:      &birthDate,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		user, err := domain.NewUser(params)
		if err != testCase.expectedError {
			t.Errorf("%s: expected error: %v, got %v", testCase.testName, testCase.expectedError, err)
		}
		if testCase.expectedError == nil {
			expected := params
			if !isEqual(user, &expected) {
				t.Errorf("%s: expected user values to match params, got mismatch", testCase.testName)
			}
		}
	}
}

func isEqual(user *domain.User, params *domain.UserParams) bool {
	if user == nil && params == nil {
		return true
	}
	if user == nil || params == nil {
		return false
	}
	dto := user.ToDTO()
	return dto.ID == params.ID &&
		dto.Identification == params.Identification &&
		dto.FirstName == params.FirstName &&
		dto.LastName == params.LastName &&
		dto.Phone == params.Phone &&
		(dto.Email == nil && params.Email == nil || dto.Email != nil && params.Email != nil && *dto.Email == *params.Email) &&
		(dto.Address == nil && params.Address == nil || dto.Address != nil && params.Address != nil && *dto.Address == *params.Address) &&
		(dto.BirthDate == nil && params.BirthDate == nil || dto.BirthDate != nil && params.BirthDate != nil && dto.BirthDate.Value.Equal(params.BirthDate.Value)) &&
		dto.CreatedAt.Equal(params.CreatedAt) &&
		dto.UpdatedAt.Equal(params.UpdatedAt)
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

	dto := user.ToDTO()
	if dto.ID != 0 {
		t.Errorf("expected generated ID to be 0, got %d", dto.ID)
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

func TestUserDTOAndFromDTO(t *testing.T) {
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

	params := domain.UserParams{
		ID:             1,
		Identification: identification,
		FirstName:      "John",
		LastName:       "Doe",
		Phone:          phone,
		Email:          &email,
		Address:        &address,
		BirthDate:      &birthDate,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	user, err := domain.NewUser(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ToDTO
	dto := user.ToDTO()
	if dto.ID != params.ID || dto.FirstName != params.FirstName {
		t.Errorf("ToDTO mismatch")
	}

	var nilUser *domain.User
	if nilUser.ToDTO() != nil {
		t.Errorf("expected nil DTO for nil User")
	}

	// FromDTO
	user2, err := domain.UserFromDTO(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !isEqual(user2, &params) {
		t.Errorf("UserFromDTO values do not match original params")
	}

	nilUser2, err := domain.UserFromDTO(nil)
	if err != nil || nilUser2 != nil {
		t.Errorf("expected nil user for nil DTO without error")
	}
}

func TestUserWith(t *testing.T) {
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

	params := domain.UserParams{
		ID:             1,
		Identification: identification,
		FirstName:      "John",
		LastName:       "Doe",
		Phone:          phone,
		Email:          &email,
		Address:        &address,
		BirthDate:      &birthDate,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	user, err := domain.NewUser(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newPhone := domain.Phone{Value: "987654321"}
	newFirstName := "Jane"

	updatedUser, err := user.With(domain.UserParams{
		Identification: identification,
		FirstName:      newFirstName,
		LastName:       "Doe",
		Phone:          newPhone,
		Email:          &email,
		Address:        &address,
		BirthDate:      &birthDate,
	})
	if err != nil {
		t.Fatalf("unexpected error on With: %v", err)
	}

	dto := updatedUser.ToDTO()
	if dto.ID != 1 {
		t.Errorf("expected ID to be kept as 1, got %d", dto.ID)
	}
	if dto.FirstName != newFirstName {
		t.Errorf("expected FirstName to be updated to %s, got %s", newFirstName, dto.FirstName)
	}
	if dto.Phone != newPhone {
		t.Errorf("expected Phone to be updated to %v, got %v", newPhone, dto.Phone)
	}
	if dto.CreatedAt != now {
		t.Errorf("expected CreatedAt to be kept as %v, got %v", now, dto.CreatedAt)
	}
	if dto.UpdatedAt.Before(now) {
		t.Errorf("expected UpdatedAt to be updated, got %v", dto.UpdatedAt)
	}
}

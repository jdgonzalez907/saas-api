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
			testName:      "success - create user",
			id:            1,
			firstName:     "John",
			lastName:      "Doe",
			expectedError: nil,
		},
		{
			testName:      "fail - invalid id less than 0",
			id:            -1,
			firstName:     "John",
			lastName:      "Doe",
			expectedError: domain.ErrInvalidUserID,
		},
		{
			testName:      "fail - empty first name",
			id:            1,
			firstName:     "",
			lastName:      "Doe",
			expectedError: domain.ErrInvalidFirstName,
		},
		{
			testName:      "fail - empty last name",
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

	testCases := []struct {
		testName      string
		firstName     string
		lastName      string
		expectedError error
	}{
		{
			testName:      "success - create user",
			firstName:     "John",
			lastName:      "Doe",
			expectedError: nil,
		},
		{
			testName:      "fail - empty first name",
			firstName:     "",
			lastName:      "Doe",
			expectedError: domain.ErrInvalidFirstName,
		},
		{
			testName:      "fail - empty last name",
			firstName:     "John",
			lastName:      "",
			expectedError: domain.ErrInvalidLastName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			user, err := domain.NewUserWithoutId(
				identification,
				tc.firstName,
				tc.lastName,
				phone,
				&email,
				&address,
				&birthDate,
			)

			if err != tc.expectedError {
				t.Errorf("expected error: %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				if user == nil {
					t.Fatal("expected user to be not nil")
				}
				dto := user.ToDTO()
				if dto.ID != 0 {
					t.Errorf("expected generated ID to be 0, got %d", dto.ID)
				}
			}
		})
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

func TestUserWithPersonalInformation(t *testing.T) {
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

	newIdentification := domain.Identification{Type: domain.IdType_CC, Number: "987654321"}
	newFirstName := "Jane"
	newLastName := "Smith"
	newAddress := domain.Address{
		Street:      "456 Main St",
		PostalCode:  nil,
		City:        "Boston",
		State:       "MA",
		Country:     "USA",
		Description: nil,
	}
	newBirthDate := domain.BirthDate{Value: time.Now().AddDate(-1, 0, 0)}

	testCases := []struct {
		testName       string
		identification domain.Identification
		firstName      string
		lastName       string
		address        *domain.Address
		birthDate      *domain.BirthDate
		expectedError  error
	}{
		{
			testName:       "success",
			identification: newIdentification,
			firstName:      newFirstName,
			lastName:       newLastName,
			address:        &newAddress,
			birthDate:      &newBirthDate,
			expectedError:  nil,
		},
		{
			testName:       "fail - invalid first name",
			identification: newIdentification,
			firstName:      "",
			lastName:       newLastName,
			address:        &newAddress,
			birthDate:      &newBirthDate,
			expectedError:  domain.ErrInvalidFirstName,
		},
		{
			testName:       "fail - invalid last name",
			identification: newIdentification,
			firstName:      newFirstName,
			lastName:       "",
			address:        &newAddress,
			birthDate:      &newBirthDate,
			expectedError:  domain.ErrInvalidLastName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			updatedUser, err := user.WithPersonalInformation(
				tc.identification,
				tc.firstName,
				tc.lastName,
				tc.address,
				tc.birthDate,
			)

			if err != tc.expectedError {
				t.Errorf("expected error: %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				dto := updatedUser.ToDTO()
				if dto.ID != 1 {
					t.Errorf("expected ID to be kept as 1, got %d", dto.ID)
				}
				if dto.FirstName != tc.firstName {
					t.Errorf("expected FirstName to be updated to %s, got %s", tc.firstName, dto.FirstName)
				}
				if dto.LastName != tc.lastName {
					t.Errorf("expected LastName to be updated to %s, got %s", tc.lastName, dto.LastName)
				}
				if dto.Identification != tc.identification {
					t.Errorf("expected Identification to be updated")
				}
				if dto.Phone != phone {
					t.Errorf("expected Phone to be kept as %v, got %v", phone, dto.Phone)
				}
				if dto.Email == nil || *dto.Email != email {
					t.Errorf("expected Email to be kept")
				}
				if dto.Address == nil || *dto.Address != *tc.address {
					t.Errorf("expected Address to be updated")
				}
				if dto.BirthDate == nil || !dto.BirthDate.Value.Equal(tc.birthDate.Value) {
					t.Errorf("expected BirthDate to be updated")
				}
				if dto.CreatedAt != now {
					t.Errorf("expected CreatedAt to be kept as %v, got %v", now, dto.CreatedAt)
				}
				if dto.UpdatedAt.Before(now) {
					t.Errorf("expected UpdatedAt to be updated, got %v", dto.UpdatedAt)
				}
			}
		})
	}
}

func TestUserWithPhone(t *testing.T) {
	identification := domain.Identification{Type: domain.IdType_CC, Number: "123456789"}
	email := domain.Email{Value: "name@domain.com"}
	phone := domain.Phone{Value: "123456789"}
	address := domain.Address{Street: "123 Main St"}
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
	updatedUser := user.WithPhone(newPhone)

	dto := updatedUser.ToDTO()
	if dto.Phone != newPhone {
		t.Errorf("expected Phone to be updated, got %v", dto.Phone)
	}
	if dto.FirstName != "John" {
		t.Errorf("expected FirstName to remain unchanged")
	}
	if dto.UpdatedAt.Before(now) {
		t.Errorf("expected UpdatedAt to be updated")
	}
}

func TestUserWithEmail(t *testing.T) {
	identification := domain.Identification{Type: domain.IdType_CC, Number: "123456789"}
	email := domain.Email{Value: "name@domain.com"}
	phone := domain.Phone{Value: "123456789"}
	address := domain.Address{Street: "123 Main St"}
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

	newEmail := domain.Email{Value: "new@domain.com"}
	updatedUser := user.WithEmail(&newEmail)

	dto := updatedUser.ToDTO()
	if dto.Email == nil || *dto.Email != newEmail {
		t.Errorf("expected Email to be updated")
	}
	if dto.FirstName != "John" {
		t.Errorf("expected FirstName to remain unchanged")
	}
	if dto.UpdatedAt.Before(now) {
		t.Errorf("expected UpdatedAt to be updated")
	}
}

package domain_test

import (
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
)

var (
	testTime          = time.Date(1995, 5, 5, 0, 0, 0, 0, time.UTC)
	identification, _ = domain.NewIdentification(domain.IdType_CC, "123456789")
	email, _          = domain.NewEmail("name@domain.com")
	phone, _          = domain.NewPhone("123456789")
	address, _        = domain.NewAddress("123 Main St", "New York", "NY", "USA", nil, nil)
	birthDate, _      = domain.NewBirthDate(testTime)
)

func TestNewUser(t *testing.T) {
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
		dto.Identification == params.Identification.ToDTO() &&
		dto.FirstName == params.FirstName &&
		dto.LastName == params.LastName &&
		dto.Phone == params.Phone.ToDTO() &&
		(dto.Email == nil && params.Email == nil || dto.Email != nil && params.Email != nil && *dto.Email == params.Email.ToDTO()) &&
		(dto.Address == nil && params.Address == nil || dto.Address != nil && params.Address != nil && *dto.Address == params.Address.ToDTO()) &&
		(dto.BirthDate == nil && params.BirthDate == nil || dto.BirthDate != nil && params.BirthDate != nil && dto.BirthDate.Value.Equal(params.BirthDate.ToDTO().Value)) &&
		dto.CreatedAt.Equal(params.CreatedAt) &&
		dto.UpdatedAt.Equal(params.UpdatedAt)
}

func TestNewUserWithoutId(t *testing.T) {
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

	// Invalid cases in UserFromDTO
	invalidCases := []struct {
		name   string
		modify func(d *domain.UserDTO)
	}{
		{
			name: "invalid identification type",
			modify: func(d *domain.UserDTO) {
				d.Identification.Type = "INVALID"
			},
		},
		{
			name: "invalid phone",
			modify: func(d *domain.UserDTO) {
				d.Phone.Value = ""
			},
		},
		{
			name: "invalid email",
			modify: func(d *domain.UserDTO) {
				d.Email = &domain.EmailDTO{Value: "invalid-email"}
			},
		},
		{
			name: "invalid address",
			modify: func(d *domain.UserDTO) {
				d.Address = &domain.AddressDTO{Street: ""}
			},
		},
		{
			name: "invalid birthdate",
			modify: func(d *domain.UserDTO) {
				d.BirthDate = &domain.BirthDateDTO{Value: time.Now()}
			},
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			validDTO := user.ToDTO()
			tc.modify(validDTO)
			_, err := domain.UserFromDTO(validDTO)
			if err == nil {
				t.Errorf("expected error for case %s, got nil", tc.name)
			}
		})
	}
}

func TestUserWithPersonalInformation(t *testing.T) {
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

	newIdentification, _ := domain.NewIdentification(domain.IdType_CC, "987654321")
	newFirstName := "Jane"
	newLastName := "Smith"
	newAddress, _ := domain.NewAddress("456 Main St", "Boston", "MA", "USA", nil, nil)
	newBirthDate, _ := domain.NewBirthDate(time.Now().AddDate(-19, 0, 0))

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
				if dto.Identification != tc.identification.ToDTO() {
					t.Errorf("expected Identification to be updated")
				}
				if dto.Phone != phone.ToDTO() {
					t.Errorf("expected Phone to be kept as %v, got %v", phone.ToDTO(), dto.Phone)
				}
				if dto.Email == nil || *dto.Email != email.ToDTO() {
					t.Errorf("expected Email to be kept")
				}
				if dto.Address == nil || *dto.Address != tc.address.ToDTO() {
					t.Errorf("expected Address to be updated")
				}
				if dto.BirthDate == nil || !dto.BirthDate.Value.Equal(tc.birthDate.ToDTO().Value) {
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

	newPhone, _ := domain.NewPhone("987654321")
	updatedUser := user.WithPhone(newPhone)

	dto := updatedUser.ToDTO()
	if dto.Phone != newPhone.ToDTO() {
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

	newEmail, _ := domain.NewEmail("new@domain.com")
	updatedUser := user.WithEmail(&newEmail)

	dto := updatedUser.ToDTO()
	if dto.Email == nil || *dto.Email != newEmail.ToDTO() {
		t.Errorf("expected Email to be updated")
	}
	if dto.FirstName != "John" {
		t.Errorf("expected FirstName to remain unchanged")
	}
	if dto.UpdatedAt.Before(now) {
		t.Errorf("expected UpdatedAt to be updated")
	}
}

func TestUserGetters(t *testing.T) {
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

	if user.ID() != params.ID {
		t.Errorf("expected ID: %d, got: %d", params.ID, user.ID())
	}
	if user.Identification() != params.Identification {
		t.Errorf("expected Identification: %v, got: %v", params.Identification, user.Identification())
	}
	if user.FirstName() != params.FirstName {
		t.Errorf("expected FirstName: %s, got: %s", params.FirstName, user.FirstName())
	}
	if user.LastName() != params.LastName {
		t.Errorf("expected LastName: %s, got: %s", params.LastName, user.LastName())
	}
	if user.Phone() != params.Phone {
		t.Errorf("expected Phone: %v, got: %v", params.Phone, user.Phone())
	}
	if user.Email() != params.Email {
		t.Errorf("expected Email: %v, got: %v", params.Email, user.Email())
	}
	if user.Address() != params.Address {
		t.Errorf("expected Address: %v, got: %v", params.Address, user.Address())
	}
	if user.BirthDate() != params.BirthDate {
		t.Errorf("expected BirthDate: %v, got: %v", params.BirthDate, user.BirthDate())
	}
	if !user.CreatedAt().Equal(params.CreatedAt) {
		t.Errorf("expected CreatedAt: %v, got: %v", params.CreatedAt, user.CreatedAt())
	}
	if !user.UpdatedAt().Equal(params.UpdatedAt) {
		t.Errorf("expected UpdatedAt: %v, got: %v", params.UpdatedAt, user.UpdatedAt())
	}
}

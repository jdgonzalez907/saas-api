package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"jdgonzalez907/saas-api/internal/users/domain"
)

var (
	identification, _ = domain.NewIdentification(domain.IDTypeCC, "123456789")
	email, _          = domain.NewEmail("name@domain.com")
	phone, _          = domain.NewPhone("57", "123456789")
	address, _        = domain.NewAddress("123 Main St", "New York", "NY", "USA", nil, nil)
	birthDate, _      = domain.NewBirthDate("1995-05-05")

	personalInfo, _ = domain.NewPersonalInformation(
		identification,
		"John",
		"Doe",
		&address,
		&birthDate,
	)
)

func TestNewUser(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		testName      string
		id            int64
		info          domain.PersonalInformation
		expectedError error
	}{
		{
			testName:      "success - create user",
			id:            1,
			info:          personalInfo,
			expectedError: nil,
		},
		{
			testName:      "fail - invalid id less than 0",
			id:            -1,
			info:          personalInfo,
			expectedError: domain.ErrInvalidUserID,
		},
		{
			testName:      "fail - invalid id equal to 0",
			id:            0,
			info:          personalInfo,
			expectedError: domain.ErrInvalidUserID,
		},
	}

	for _, testCase := range testCases {
		user, err := domain.NewUser(
			testCase.id,
			testCase.info,
			phone,
			&email,
			now,
			now,
		)
		if err != testCase.expectedError {
			t.Errorf("%s: expected error: %v, got %v", testCase.testName, testCase.expectedError, err)
		}
		if testCase.expectedError == nil {
			if user == nil {
				t.Errorf("%s: expected user to be not nil", testCase.testName)
				continue
			}
			dto := user.ToDTO()
			if dto.ID != testCase.id {
				t.Errorf("%s: expected ID %d, got %d", testCase.testName, testCase.id, dto.ID)
			}
		}
	}
}

func TestNewUserWithoutId(t *testing.T) {
	user, err := domain.NewUserWithoutID(
		personalInfo,
		phone,
		&email,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be not nil")
	}

	dto := user.ToDTO()
	if dto.ID != 0 {
		t.Errorf("expected generated ID to be 0, got %d", dto.ID)
	}
}

func TestAssignID(t *testing.T) {
	user, err := domain.NewUserWithoutID(personalInfo, phone, &email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID() != 0 {
		t.Fatalf("expected initial ID to be 0, got %d", user.ID())
	}

	user.AssignID(42)

	if user.ID() != 42 {
		t.Errorf("expected ID to be 42 after AssignID, got %d", user.ID())
	}
}

func TestUserDTOAndFromDTO(t *testing.T) {
	now := time.Now()

	user, err := domain.NewUser(
		1,
		personalInfo,
		phone,
		&email,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ToDTO
	dto := user.ToDTO()
	if dto.ID != 1 || dto.FirstName != personalInfo.ToDTO().FirstName {
		t.Errorf("ToDTO mismatch")
	}

	var nilUser *domain.User
	if nilUser.ToDTO() != nil {
		t.Errorf("expected nil DTO for nil User")
	}

	// FromDTO
	restoredUser, err := domain.UserFromDTO(dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if restoredUser == nil {
		t.Fatal("expected restored user to be not nil")
	}

	restoredDTO := restoredUser.ToDTO()
	if restoredDTO.ID != dto.ID ||
		restoredDTO.FirstName != dto.FirstName ||
		restoredDTO.LastName != dto.LastName ||
		restoredDTO.Identification != dto.Identification ||
		restoredDTO.Phone != dto.Phone ||
		(restoredDTO.Email == nil) != (dto.Email == nil) ||
		(restoredDTO.Email != nil && dto.Email != nil && *restoredDTO.Email != *dto.Email) ||
		!restoredDTO.CreatedAt.Equal(dto.CreatedAt) ||
		!restoredDTO.UpdatedAt.Equal(dto.UpdatedAt) {
		t.Errorf("UserFromDTO values do not match original")
	}

	nilUserFromDTO, err := domain.UserFromDTO(nil)
	if err != nil || nilUserFromDTO != nil {
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
				d.Phone.Number = ""
			},
		},
		{
			name: "invalid email",
			modify: func(d *domain.UserDTO) {
				badEmail := domain.EmailDTO("invalid-email")
				d.Email = &badEmail
			},
		},
		{
			name: "invalid address",
			modify: func(d *domain.UserDTO) {
				d.Address = &domain.AddressDTO{Street: ""}
			},
		},
		{
			name: "invalid birthdate format",
			modify: func(d *domain.UserDTO) {
				bd := domain.BirthDateDTO("invalid-date")
				d.BirthDate = &bd
			},
		},
		{
			name: "invalid birthdate underage",
			modify: func(d *domain.UserDTO) {
				bd := domain.BirthDateDTO(time.Now().Format("2006-01-02"))
				d.BirthDate = &bd
			},
		},
		{
			name: "invalid first name",
			modify: func(d *domain.UserDTO) {
				d.FirstName = ""
			},
		},
		{
			name: "invalid last name",
			modify: func(d *domain.UserDTO) {
				d.LastName = ""
			},
		},
		{
			name: "invalid user id",
			modify: func(d *domain.UserDTO) {
				d.ID = -1
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

func TestUpdatePersonalInformation(t *testing.T) {
	now := time.Now()

	user, err := domain.NewUser(
		1,
		personalInfo,
		phone,
		&email,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newIdentification, _ := domain.NewIdentification(domain.IDTypeCC, "987654321")
	newFirstName := "Jane"
	newLastName := "Smith"
	newAddress, _ := domain.NewAddress("456 Main St", "Boston", "MA", "USA", nil, nil)
	newBirthDate, _ := domain.NewBirthDate(time.Now().AddDate(-19, 0, 0).Format("2006-01-02"))

	validInfo, _ := domain.NewPersonalInformation(
		newIdentification,
		newFirstName,
		newLastName,
		&newAddress,
		&newBirthDate,
	)

	validInfoWithNil, _ := domain.NewPersonalInformation(
		newIdentification,
		newFirstName,
		newLastName,
		nil,
		nil,
	)

	testCases := []struct {
		testName        string
		info            domain.PersonalInformation
		expectedAddress *domain.Address
		expectedBD      *domain.BirthDate
	}{
		{
			testName:        "success - with valid personal information",
			info:            validInfo,
			expectedAddress: &newAddress,
			expectedBD:      &newBirthDate,
		},
		{
			testName:        "success - with nil fields",
			info:            validInfoWithNil,
			expectedAddress: nil,
			expectedBD:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			updatedUser := user.UpdatePersonalInformation(tc.info)
			dto := updatedUser.ToDTO()

			if dto.ID != 1 {
				t.Errorf("expected ID to be kept as 1, got %d", dto.ID)
			}
			if dto.FirstName != newFirstName {
				t.Errorf("expected FirstName to be updated to %s, got %s", newFirstName, dto.FirstName)
			}
			if dto.LastName != newLastName {
				t.Errorf("expected LastName to be updated to %s, got %s", newLastName, dto.LastName)
			}
			if dto.Identification != newIdentification.ToDTO() {
				t.Errorf("expected Identification to be updated")
			}
			if dto.Phone != phone.ToDTO() {
				t.Errorf("expected Phone to be kept as %v, got %v", phone.ToDTO(), dto.Phone)
			}
			if dto.Email == nil || *dto.Email != email.ToDTO() {
				t.Errorf("expected Email to be kept")
			}
			if tc.expectedAddress == nil {
				if dto.Address != nil {
					t.Errorf("expected nil Address, got %v", dto.Address)
				}
			} else {
				if dto.Address == nil || *dto.Address != tc.expectedAddress.ToDTO() {
					t.Errorf("expected Address to be updated")
				}
			}
			if tc.expectedBD == nil {
				if dto.BirthDate != nil {
					t.Errorf("expected nil BirthDate, got %v", dto.BirthDate)
				}
			} else {
				if dto.BirthDate == nil || *dto.BirthDate != tc.expectedBD.ToDTO() {
					t.Errorf("expected BirthDate to be updated")
				}
			}
		})
	}
}

func TestChangePhone(t *testing.T) {
	now := time.Now()

	user, err := domain.NewUser(
		1,
		personalInfo,
		phone,
		&email,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newPhone, _ := domain.NewPhone("57", "987654321")
	updatedUser := user.ChangePhone(newPhone)

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

func TestChangeEmail(t *testing.T) {
	now := time.Now()

	user, err := domain.NewUser(
		1,
		personalInfo,
		phone,
		&email,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newEmail, _ := domain.NewEmail("new@domain.com")
	updatedUser := user.ChangeEmail(&newEmail)

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

	user, err := domain.NewUser(
		1,
		personalInfo,
		phone,
		&email,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	infoDTO := user.PersonalInformation().ToDTO()
	expectedInfoDTO := personalInfo.ToDTO()

	if user.ID() != 1 {
		t.Errorf("expected ID: %d, got: %d", 1, user.ID())
	}
	if user.Identification().ToDTO() != expectedInfoDTO.Identification {
		t.Errorf("expected Identification")
	}
	if user.FirstName() != expectedInfoDTO.FirstName {
		t.Errorf("expected FirstName")
	}
	if user.LastName() != expectedInfoDTO.LastName {
		t.Errorf("expected LastName")
	}
	if user.Phone() != phone {
		t.Errorf("expected Phone")
	}
	if user.Email() != &email {
		t.Errorf("expected Email")
	}
	if user.Address() != nil && expectedInfoDTO.Address != nil && user.Address().ToDTO() != *expectedInfoDTO.Address {
		t.Errorf("expected Address")
	}
	if user.BirthDate() != nil && expectedInfoDTO.BirthDate != nil && user.BirthDate().ToDTO() != *expectedInfoDTO.BirthDate {
		t.Errorf("expected BirthDate")
	}
	if !user.CreatedAt().Equal(now) {
		t.Errorf("expected CreatedAt")
	}
	if !user.UpdatedAt().Equal(now) {
		t.Errorf("expected UpdatedAt")
	}

	if infoDTO.FirstName != expectedInfoDTO.FirstName {
		t.Errorf("expected FirstName %s, got %s", expectedInfoDTO.FirstName, infoDTO.FirstName)
	}
	if infoDTO.LastName != expectedInfoDTO.LastName {
		t.Errorf("expected LastName %s, got %s", expectedInfoDTO.LastName, infoDTO.LastName)
	}
	if infoDTO.Identification != expectedInfoDTO.Identification {
		t.Errorf("expected Identification %v, got %v", expectedInfoDTO.Identification, infoDTO.Identification)
	}
	if (infoDTO.Address == nil) != (expectedInfoDTO.Address == nil) {
		t.Errorf("expected Address nil status to match")
	} else if infoDTO.Address != nil {
		if *infoDTO.Address != *expectedInfoDTO.Address {
			t.Errorf("expected Address value to match")
		}
	}
	if (infoDTO.BirthDate == nil) != (expectedInfoDTO.BirthDate == nil) {
		t.Errorf("expected BirthDate nil status to match")
	} else if infoDTO.BirthDate != nil {
		if *infoDTO.BirthDate != *expectedInfoDTO.BirthDate {
			t.Errorf("expected BirthDate value to match")
		}
	}

	userFullName := user.FullName()
	assert.Equal(t, user.FullName(), userFullName)
}

func TestValidateAssignedUserID(t *testing.T) {
	testCases := []struct {
		name          string
		id            int64
		expectedError error
	}{
		{
			name:          "success - valid id",
			id:            1,
			expectedError: nil,
		},
		{
			name:          "fail - unassigned id",
			id:            0,
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name:          "fail - negative id",
			id:            -1,
			expectedError: domain.ErrInvalidUserID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := domain.ValidateAssignedUserID(tc.id)
			if err != tc.expectedError {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}

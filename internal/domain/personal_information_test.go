package domain_test

import (
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
)

func TestPersonalInformation(t *testing.T) {
	identification, _ := domain.NewIdentification(domain.IdType_CC, "123456")
	address, _ := domain.NewAddress("St", "City", "State", "Country", nil, nil)
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-25, 0, 0).Format("2006-01-02"))

	testCases := []struct {
		testName      string
		firstName     string
		lastName      string
		address       *domain.Address
		birthDate     *domain.BirthDate
		expectedError error
	}{
		{
			testName:      "success - create personal info and map to DTO",
			firstName:     "John",
			lastName:      "Doe",
			address:       &address,
			birthDate:     &birthDate,
			expectedError: nil,
		},
		{
			testName:      "success - create personal info with nil fields",
			firstName:     "John",
			lastName:      "Doe",
			address:       nil,
			birthDate:     nil,
			expectedError: nil,
		},
		{
			testName:      "fail - empty first name",
			firstName:     "",
			lastName:      "Doe",
			address:       &address,
			birthDate:     &birthDate,
			expectedError: domain.ErrInvalidFirstName,
		},
		{
			testName:      "fail - empty last name",
			firstName:     "John",
			lastName:      "",
			address:       &address,
			birthDate:     &birthDate,
			expectedError: domain.ErrInvalidLastName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			info, err := domain.NewPersonalInformation(
				identification,
				tc.firstName,
				tc.lastName,
				tc.address,
				tc.birthDate,
			)
			if err != tc.expectedError {
				t.Fatalf("expected error %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				dto := info.ToDTO()
				if dto.FirstName != tc.firstName {
					t.Errorf("expected FirstName %s, got %s", tc.firstName, dto.FirstName)
				}
				if dto.LastName != tc.lastName {
					t.Errorf("expected LastName %s, got %s", tc.lastName, dto.LastName)
				}
				if tc.address == nil {
					if dto.Address != nil {
						t.Errorf("expected nil Address, got %+v", dto.Address)
					}
				} else {
					if dto.Address == nil || dto.Address.Street != tc.address.ToDTO().Street {
						t.Errorf("expected Address %s, got %v", tc.address.ToDTO().Street, dto.Address)
					}
				}
				if tc.birthDate == nil {
					if dto.BirthDate != nil {
						t.Errorf("expected nil BirthDate, got %+v", dto.BirthDate)
					}
				} else {
					if dto.BirthDate == nil {
						t.Errorf("expected non-nil BirthDate")
					}
				}
			}
		})
	}
}

package domain_test

import (
	"testing"
	"time"

	"jdgonzalez907/saas-api/internal/users/domain"
)

func TestPersonalInformation(t *testing.T) {
	identification, _ := domain.NewIdentification(domain.IDTypeCC, "123456")
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

func TestPersonalInformation_Equals(t *testing.T) {
	id1, _ := domain.NewIdentification(domain.IDTypeCC, "123456")
	id2, _ := domain.NewIdentification(domain.IDTypePASSPORT, "123456")
	addr1, _ := domain.NewAddress("St1", "City", "State", "Country", nil, nil)
	addr2, _ := domain.NewAddress("St2", "City", "State", "Country", nil, nil)
	bd1, _ := domain.NewBirthDate("2000-01-01")
	bd2, _ := domain.NewBirthDate("1999-01-01")

	piBase, _ := domain.NewPersonalInformation(id1, "John", "Doe", &addr1, &bd1)
	piSame, _ := domain.NewPersonalInformation(id1, "John", "Doe", &addr1, &bd1)
	piDiffName, _ := domain.NewPersonalInformation(id1, "Jane", "Doe", &addr1, &bd1)
	piDiffID, _ := domain.NewPersonalInformation(id2, "John", "Doe", &addr1, &bd1)
	piNilAddr, _ := domain.NewPersonalInformation(id1, "John", "Doe", nil, &bd1)
	piDiffAddr, _ := domain.NewPersonalInformation(id1, "John", "Doe", &addr2, &bd1)
	piNilBD, _ := domain.NewPersonalInformation(id1, "John", "Doe", &addr1, nil)
	piDiffBD, _ := domain.NewPersonalInformation(id1, "John", "Doe", &addr1, &bd2)

	testCases := []struct {
		testName string
		pi1      domain.PersonalInformation
		pi2      domain.PersonalInformation
		expected bool
	}{
		{
			testName: "success - identical personal info",
			pi1:      piBase,
			pi2:      piSame,
			expected: true,
		},
		{
			testName: "fail - different first name",
			pi1:      piBase,
			pi2:      piDiffName,
			expected: false,
		},
		{
			testName: "fail - different identification",
			pi1:      piBase,
			pi2:      piDiffID,
			expected: false,
		},
		{
			testName: "fail - nil address vs non-nil",
			pi1:      piBase,
			pi2:      piNilAddr,
			expected: false,
		},
		{
			testName: "fail - different address",
			pi1:      piBase,
			pi2:      piDiffAddr,
			expected: false,
		},
		{
			testName: "fail - nil birth date vs non-nil",
			pi1:      piBase,
			pi2:      piNilBD,
			expected: false,
		},
		{
			testName: "fail - different birth date",
			pi1:      piBase,
			pi2:      piDiffBD,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result := tc.pi1.Equals(tc.pi2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

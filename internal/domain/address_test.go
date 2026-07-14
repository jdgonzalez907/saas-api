package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewAddress(t *testing.T) {
	postalCodeVal := "postal_code"
	descriptionVal := "description"
	emptyStrVal := ""

	pPostalCode := &postalCodeVal
	pDescription := &descriptionVal
	pEmptyStr := &emptyStrVal

	testCases := []struct {
		testName      string
		street        string
		city          string
		state         string
		country       string
		postalCode    *string
		description   *string
		expectedError error
	}{
		{
			testName:      "success - create address",
			street:        "street",
			city:          "city",
			state:         "state",
			country:       "country",
			postalCode:    nil,
			description:   nil,
			expectedError: nil,
		},
		{
			testName:      "success - with postal code and description",
			street:        "street",
			city:          "city",
			state:         "state",
			country:       "country",
			postalCode:    pPostalCode,
			description:   pDescription,
			expectedError: nil,
		},
		{
			testName:      "fail - empty street",
			street:        "",
			city:          "city",
			state:         "state",
			country:       "country",
			postalCode:    pPostalCode,
			description:   pDescription,
			expectedError: domain.ErrInvalidStreet,
		},
		{
			testName:      "fail - empty city",
			street:        "street",
			city:          "",
			state:         "state",
			country:       "country",
			postalCode:    pPostalCode,
			description:   pDescription,
			expectedError: domain.ErrInvalidCity,
		},
		{
			testName:      "fail - empty state",
			street:        "street",
			city:          "city",
			state:         "",
			country:       "country",
			postalCode:    pPostalCode,
			description:   pDescription,
			expectedError: domain.ErrInvalidState,
		},
		{
			testName:      "fail - empty country",
			street:        "street",
			city:          "city",
			state:         "state",
			country:       "",
			postalCode:    pPostalCode,
			description:   pDescription,
			expectedError: domain.ErrInvalidCountry,
		},
		{
			testName:      "fail - empty postal code pointer",
			street:        "street",
			city:          "city",
			state:         "state",
			country:       "country",
			postalCode:    pEmptyStr,
			description:   pDescription,
			expectedError: domain.ErrInvalidPostalCode,
		},
		{
			testName:      "fail - empty description pointer",
			street:        "street",
			city:          "city",
			state:         "state",
			country:       "country",
			postalCode:    pPostalCode,
			description:   pEmptyStr,
			expectedError: domain.ErrInvalidDescription,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			address, err := domain.NewAddress(tc.street, tc.city, tc.state, tc.country, tc.postalCode, tc.description)
			if err != tc.expectedError {
				t.Fatalf("expected error: %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				dto := address.ToDTO()
				if dto.Street != tc.street {
					t.Errorf("expected street: %s, got: %s", tc.street, dto.Street)
				}
				if dto.City != tc.city {
					t.Errorf("expected city: %s, got: %s", tc.city, dto.City)
				}
				if dto.State != tc.state {
					t.Errorf("expected state: %s, got: %s", tc.state, dto.State)
				}
				if dto.Country != tc.country {
					t.Errorf("expected country: %s, got: %s", tc.country, dto.Country)
				}
				if dto.PostalCode != tc.postalCode {
					t.Errorf("expected postal code: %v, got: %v", tc.postalCode, dto.PostalCode)
				}
				if dto.Description != tc.description {
					t.Errorf("expected description: %v, got: %v", tc.description, dto.Description)
				}
			}
		})
	}
}

func TestAddress_Equals(t *testing.T) {
	pc1 := "12345"
	pc2 := "12345"
	pc3 := "54321"
	desc1 := "house"
	desc2 := "house"
	desc3 := "work"

	addrBase, _ := domain.NewAddress("St 1", "City", "State", "US", &pc1, &desc1)
	addrSame, _ := domain.NewAddress("St 1", "City", "State", "US", &pc2, &desc2)
	addrDiffStreet, _ := domain.NewAddress("St 2", "City", "State", "US", &pc1, &desc1)
	addrNilPC, _ := domain.NewAddress("St 1", "City", "State", "US", nil, &desc1)
	addrDiffPC, _ := domain.NewAddress("St 1", "City", "State", "US", &pc3, &desc1)
	addrNilDesc, _ := domain.NewAddress("St 1", "City", "State", "US", &pc1, nil)
	addrDiffDesc, _ := domain.NewAddress("St 1", "City", "State", "US", &pc1, &desc3)

	testCases := []struct {
		testName string
		addr1    domain.Address
		addr2    domain.Address
		expected bool
	}{
		{
			testName: "success - identical addresses",
			addr1:    addrBase,
			addr2:    addrSame,
			expected: true,
		},
		{
			testName: "fail - different street",
			addr1:    addrBase,
			addr2:    addrDiffStreet,
			expected: false,
		},
		{
			testName: "fail - nil postal code vs non-nil",
			addr1:    addrBase,
			addr2:    addrNilPC,
			expected: false,
		},
		{
			testName: "fail - different postal code",
			addr1:    addrBase,
			addr2:    addrDiffPC,
			expected: false,
		},
		{
			testName: "fail - nil description vs non-nil",
			addr1:    addrBase,
			addr2:    addrNilDesc,
			expected: false,
		},
		{
			testName: "fail - different description",
			addr1:    addrBase,
			addr2:    addrDiffDesc,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result := tc.addr1.Equals(tc.addr2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

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
		testName       string
		input          domain.Address
		expectedError  error
		expectedOutput domain.Address
	}{
		{
			testName: "create address",
			input: domain.Address{
				Street:      "street",
				City:        "city",
				State:       "state",
				Country:     "country",
				PostalCode:  nil,
				Description: nil,
			},
			expectedError: nil,
			expectedOutput: domain.Address{
				Street:      "street",
				City:        "city",
				State:       "state",
				Country:     "country",
				PostalCode:  nil,
				Description: nil,
			},
		},
		{
			testName: "create address with postal code and description",
			input: domain.Address{
				Street:      "street",
				City:        "city",
				State:       "state",
				Country:     "country",
				PostalCode:  pPostalCode,
				Description: pDescription,
			},
			expectedError: nil,
			expectedOutput: domain.Address{
				Street:      "street",
				City:        "city",
				State:       "state",
				Country:     "country",
				PostalCode:  pPostalCode,
				Description: pDescription,
			},
		},
		{
			testName: "fail to create address with empty street",
			input: domain.Address{
				Street:      "",
				City:        "city",
				State:       "state",
				Country:     "country",
				PostalCode:  pPostalCode,
				Description: pDescription,
			},
			expectedError:  domain.ErrInvalidStreet,
			expectedOutput: domain.Address{},
		},
		{
			testName: "fail to create address with empty city",
			input: domain.Address{
				Street:      "street",
				City:        "",
				State:       "state",
				Country:     "country",
				PostalCode:  pPostalCode,
				Description: pDescription,
			},
			expectedError:  domain.ErrInvalidCity,
			expectedOutput: domain.Address{},
		},
		{
			testName: "fail to create address with empty state",
			input: domain.Address{
				Street:      "street",
				City:        "city",
				State:       "",
				Country:     "country",
				PostalCode:  pPostalCode,
				Description: pDescription,
			},
			expectedError:  domain.ErrInvalidState,
			expectedOutput: domain.Address{},
		},
		{
			testName: "fail to create address with empty country",
			input: domain.Address{
				Street:      "street",
				City:        "city",
				State:       "state",
				Country:     "",
				PostalCode:  pPostalCode,
				Description: pDescription,
			},
			expectedError:  domain.ErrInvalidCountry,
			expectedOutput: domain.Address{},
		},
		{
			testName: "fail to create address with empty postal code",
			input: domain.Address{
				Street:      "street",
				City:        "city",
				State:       "state",
				Country:     "country",
				PostalCode:  pEmptyStr,
				Description: pDescription,
			},
			expectedError:  domain.ErrInvalidPostalCode,
			expectedOutput: domain.Address{},
		},
		{
			testName: "fail to create address with empty description",
			input: domain.Address{
				Street:      "street",
				City:        "city",
				State:       "state",
				Country:     "country",
				PostalCode:  pPostalCode,
				Description: pEmptyStr,
			},
			expectedError:  domain.ErrInvalidDescription,
			expectedOutput: domain.Address{},
		},
	}

	for _, testCase := range testCases {
		address, err := domain.NewAddress(testCase.input.Street, testCase.input.City, testCase.input.State, testCase.input.Country, testCase.input.PostalCode, testCase.input.Description)
		if err != testCase.expectedError {
			t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
		}
		if address != testCase.expectedOutput {
			t.Errorf("expected address: %v, got %v", testCase.expectedOutput, address)
		}
	}
}

func TestAddressString(t *testing.T) {
	postalCodeVal := "12345"
	descriptionVal := "Apt 4B"

	testCases := []struct {
		testName string
		input    domain.Address
		expected string
	}{
		{
			testName: "address with street city state country",
			input: domain.Address{
				Street:  "123 Main St",
				City:    "Springfield",
				State:   "IL",
				Country: "USA",
			},
			expected: "123 Main St, Springfield, IL, USA",
		},
		{
			testName: "address with all fields",
			input: domain.Address{
				Street:      "123 Main St",
				City:        "Springfield",
				State:       "IL",
				Country:     "USA",
				PostalCode:  &postalCodeVal,
				Description: &descriptionVal,
			},
			expected: "123 Main St, Springfield, IL, USA Apt 4B 12345",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result := tc.input.String()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

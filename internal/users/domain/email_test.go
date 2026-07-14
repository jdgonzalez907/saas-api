package domain_test

import (
	"testing"

	"jdgonzalez907/saas-api/internal/users/domain"
)

func TestNewEmail(t *testing.T) {
	testCases := []struct {
		testName      string
		input         string
		expectedError error
	}{
		{
			testName:      "success - create email",
			input:         "name@domain.com",
			expectedError: nil,
		},
		{
			testName:      "fail - empty value",
			input:         "",
			expectedError: domain.ErrInvalidEmail,
		},
		{
			testName:      "fail - invalid format",
			input:         "invalid-email-format",
			expectedError: domain.ErrInvalidEmail,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			email, err := domain.NewEmail(tc.input)
			if err != tc.expectedError {
				t.Fatalf("expected error: %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				dto := email.ToDTO()
				if string(dto) != tc.input {
					t.Errorf("expected DTO value: %s, got: %s", tc.input, dto)
				}
				if email.Value() != tc.input {
					t.Errorf("expected Value(): %s, got: %s", tc.input, email.Value())
				}
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	emailBase, _ := domain.NewEmail("test@example.com")
	emailSame, _ := domain.NewEmail("test@example.com")
	emailDiff, _ := domain.NewEmail("other@example.com")

	testCases := []struct {
		testName string
		email1   domain.Email
		email2   domain.Email
		expected bool
	}{
		{
			testName: "success - identical emails",
			email1:   emailBase,
			email2:   emailSame,
			expected: true,
		},
		{
			testName: "fail - different email value",
			email1:   emailBase,
			email2:   emailDiff,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result := tc.email1.Equals(tc.email2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

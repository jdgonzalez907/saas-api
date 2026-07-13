package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
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

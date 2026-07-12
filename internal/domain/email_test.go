package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewEmail(t *testing.T) {
	testCases := []struct {
		testName       string
		input          string
		expectedError  error
		expectedOutput domain.Email
	}{
		{
			testName:       "create email",
			input:          "name@domain.com",
			expectedError:  nil,
			expectedOutput: domain.Email{Value: "name@domain.com"},
		},
		{
			testName:       "fail to create email with empty value",
			input:          "",
			expectedError:  domain.ErrInvalidEmail,
			expectedOutput: domain.Email{},
		},
		{
			testName:       "fail to create email with invalid format",
			input:          "invalid-email-format",
			expectedError:  domain.ErrInvalidEmail,
			expectedOutput: domain.Email{},
		},
	}

	for _, testCase := range testCases {
		email, err := domain.NewEmail(testCase.input)
		if err != testCase.expectedError {
			t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
		}
		if email != testCase.expectedOutput {
			t.Errorf("expected email: %v, got %v", testCase.expectedOutput, email)
		}
	}
}

func TestEmailString(t *testing.T) {
	email := domain.Email{Value: "test@domain.com"}
	expected := "test@domain.com"
	result := email.String()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

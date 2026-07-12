package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewPhone(t *testing.T) {
	testCases := []struct {
		testName       string
		input          string
		expectedError  error
		expectedOutput domain.Phone
	}{
		{
			testName:       "create phone",
			input:          "123456789",
			expectedError:  nil,
			expectedOutput: domain.Phone{Value: "123456789"},
		},
		{
			testName:       "fail to create phone with empty value",
			input:          "",
			expectedError:  domain.ErrInvalidPhone,
			expectedOutput: domain.Phone{},
		},
	}

	for _, testCase := range testCases {
		phone, err := domain.NewPhone(testCase.input)
		if err != testCase.expectedError {
			t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
		}
		if phone != testCase.expectedOutput {
			t.Errorf("expected phone: %v, got %v", testCase.expectedOutput, phone)
		}
	}
}

func TestPhoneString(t *testing.T) {
	phone := domain.Phone{Value: "12345"}
	expected := "12345"
	result := phone.String()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

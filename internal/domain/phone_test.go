package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewPhone(t *testing.T) {
	testCases := []struct {
		testName      string
		input         string
		expectedError error
	}{
		{
			testName:      "success - create phone",
			input:         "123456789",
			expectedError: nil,
		},
		{
			testName:      "fail - empty value",
			input:         "",
			expectedError: domain.ErrInvalidPhone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			phone, err := domain.NewPhone(tc.input)
			if err != tc.expectedError {
				t.Fatalf("expected error: %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				dto := phone.ToDTO()
				if dto.Value != tc.input {
					t.Errorf("expected DTO value: %s, got: %s", tc.input, dto.Value)
				}
			}
		})
	}
}

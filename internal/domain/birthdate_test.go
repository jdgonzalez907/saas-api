package domain_test

import (
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewBirthDate(t *testing.T) {
	today := time.Now()
	date18yearsAgo := today.AddDate(-18, 0, -1)
	date17yearsAgo := today.AddDate(-18, 0, 0)

	testCases := []struct {
		testName      string
		input         time.Time
		expectedError error
	}{
		{
			testName:      "success - create birth date",
			input:         date18yearsAgo,
			expectedError: nil,
		},
		{
			testName:      "fail - age less than 18",
			input:         date17yearsAgo,
			expectedError: domain.ErrInvalidBirthDate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			birthDate, err := domain.NewBirthDate(tc.input)
			if err != tc.expectedError {
				t.Fatalf("expected error: %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				dto := birthDate.ToDTO()
				if !dto.Value.Equal(tc.input) {
					t.Errorf("expected DTO value: %v, got %v", tc.input, dto.Value)
				}
			}
		})
	}
}

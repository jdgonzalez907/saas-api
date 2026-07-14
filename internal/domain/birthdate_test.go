package domain_test

import (
	"testing"
	"time"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewBirthDate(t *testing.T) {
	today := time.Now().UTC()
	date18yearsAgoStr := today.AddDate(-18, 0, -1).Format("2006-01-02")
	date17yearsAgoStr := today.AddDate(-18, 0, 1).Format("2006-01-02")

	testCases := []struct {
		testName      string
		input         string
		expectedError error
	}{
		{
			testName:      "success - create birth date",
			input:         date18yearsAgoStr,
			expectedError: nil,
		},
		{
			testName:      "fail - age less than 18",
			input:         date17yearsAgoStr,
			expectedError: domain.ErrUserUnderage,
		},
		{
			testName:      "fail - invalid format",
			input:         "invalid-date",
			expectedError: domain.ErrInvalidBirthDateFormat,
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
				if string(dto) != tc.input {
					t.Errorf("expected DTO value: %v, got %v", tc.input, dto)
				}
				if birthDate.Formatted() != tc.input {
					t.Errorf("expected Formatted(): %s, got: %s", tc.input, birthDate.Formatted())
				}
				if birthDate.Time().IsZero() {
					t.Error("expected Time() to be non-zero")
				}
			}
		})
	}
}

func TestBirthDate_Equals(t *testing.T) {
	bdBase, _ := domain.NewBirthDate("2000-01-01")
	bdSame, _ := domain.NewBirthDate("2000-01-01")
	bdDiff, _ := domain.NewBirthDate("1999-12-31")

	testCases := []struct {
		testName string
		bd1      domain.BirthDate
		bd2      domain.BirthDate
		expected bool
	}{
		{
			testName: "success - identical birthdates",
			bd1:      bdBase,
			bd2:      bdSame,
			expected: true,
		},
		{
			testName: "fail - different date",
			bd1:      bdBase,
			bd2:      bdDiff,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result := tc.bd1.Equals(tc.bd2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestNewBirthDate_ExactAge(t *testing.T) {
	today := time.Now().UTC()
	birthdateStr := today.AddDate(-18, 0, 0).Format("2006-01-02")
	bd, err := domain.NewBirthDate(birthdateStr)
	if err != nil {
		t.Fatalf("expected no error for exactly 18 years ago, got %v", err)
	}
	if bd.Formatted() != birthdateStr {
		t.Errorf("expected %s, got %s", birthdateStr, bd.Formatted())
	}
}


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
		testName       string
		input          time.Time
		expectedError  error
		expectedOutput domain.BirthDate
	}{
		{
			testName:       "create birth date",
			input:          date18yearsAgo,
			expectedError:  nil,
			expectedOutput: domain.BirthDate{Value: date18yearsAgo},
		},
		{
			testName:       "fail to create birth date with age less than 18",
			input:          date17yearsAgo,
			expectedError:  domain.ErrInvalidBirthDate,
			expectedOutput: domain.BirthDate{},
		},
	}

	for _, testCase := range testCases {
		birthDate, err := domain.NewBirthDate(testCase.input)
		if err != testCase.expectedError {
			t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
		}
		if birthDate != testCase.expectedOutput {
			t.Errorf("expected birth date: %v, got %v", testCase.expectedOutput, birthDate)
		}
	}
}

func TestBirthDateString(t *testing.T) {
	timeVal := time.Date(1990, 10, 15, 0, 0, 0, 0, time.UTC)
	birthDate := domain.BirthDate{Value: timeVal}
	expected := "1990-10-15"
	result := birthDate.String()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

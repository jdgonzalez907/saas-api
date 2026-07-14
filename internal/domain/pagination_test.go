package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewPagination(t *testing.T) {
	lastID := int64(42)
	lastIDZero := int64(0)
	lastIDNeg := int64(-10)
	limitVal10 := int32(10)
	limitVal25 := int32(25)
	limitVal50 := int32(50)
	limitVal15 := int32(15)
	limitValNeg5 := int32(-5)

	testCases := []struct {
		testName      string
		inputLastID   *int64
		inputLimit    *int32
		expectedLimit int32
		expectedError error
	}{
		{
			testName:      "success - limit is 10",
			inputLastID:   &lastID,
			inputLimit:    &limitVal10,
			expectedLimit: 10,
			expectedError: nil,
		},
		{
			testName:      "success - limit is 25",
			inputLastID:   &lastID,
			inputLimit:    &limitVal25,
			expectedLimit: 25,
			expectedError: nil,
		},
		{
			testName:      "success - limit is 50",
			inputLastID:   &lastID,
			inputLimit:    &limitVal50,
			expectedLimit: 50,
			expectedError: nil,
		},
		{
			testName:      "fail - invalid limit returns error",
			inputLastID:   &lastID,
			inputLimit:    &limitVal15,
			expectedLimit: 0,
			expectedError: domain.ErrInvalidPaginationLimit,
		},
		{
			testName:      "fail - negative limit returns error",
			inputLastID:   &lastID,
			inputLimit:    &limitValNeg5,
			expectedLimit: 0,
			expectedError: domain.ErrInvalidPaginationLimit,
		},
		{
			testName:      "success - nil limit defaults to 10",
			inputLastID:   &lastID,
			inputLimit:    nil,
			expectedLimit: 10,
			expectedError: nil,
		},
		{
			testName:      "success - nil lastID is accepted",
			inputLastID:   nil,
			inputLimit:    &limitVal25,
			expectedLimit: 25,
			expectedError: nil,
		},
		{
			testName:      "fail - lastID zero returns error",
			inputLastID:   &lastIDZero,
			inputLimit:    &limitVal25,
			expectedLimit: 0,
			expectedError: domain.ErrInvalidPaginationCursor,
		},
		{
			testName:      "fail - lastID negative returns error",
			inputLastID:   &lastIDNeg,
			inputLimit:    &limitVal25,
			expectedLimit: 0,
			expectedError: domain.ErrInvalidPaginationCursor,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			pagination, err := domain.NewPagination(tc.inputLastID, tc.inputLimit)

			if err != tc.expectedError {
				t.Fatalf("expected error: %v, got: %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				if pagination.Limit() != tc.expectedLimit {
					t.Errorf("expected limit: %d, got: %d", tc.expectedLimit, pagination.Limit())
				}

				if tc.inputLastID == nil {
					if pagination.LastID() != nil {
						t.Errorf("expected lastID to be nil, got: %v", pagination.LastID())
					}
				} else {
					if pagination.LastID() == nil || *pagination.LastID() != *tc.inputLastID {
						t.Errorf("expected lastID: %d, got: %v", *tc.inputLastID, pagination.LastID())
					}
				}
			}
		})
	}
}

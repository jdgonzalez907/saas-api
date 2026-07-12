package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewPagination(t *testing.T) {
	lastID := 42

	testCases := []struct {
		testName      string
		inputLastID   *int
		inputLimit    int
		expectedLimit int
	}{
		{
			testName:      "success - limit is 10",
			inputLastID:   &lastID,
			inputLimit:    10,
			expectedLimit: 10,
		},
		{
			testName:      "success - limit is 25",
			inputLastID:   &lastID,
			inputLimit:    25,
			expectedLimit: 25,
		},
		{
			testName:      "success - limit is 50",
			inputLastID:   &lastID,
			inputLimit:    50,
			expectedLimit: 50,
		},
		{
			testName:      "success - invalid limit defaults to 10",
			inputLastID:   &lastID,
			inputLimit:    15,
			expectedLimit: 10,
		},
		{
			testName:      "success - negative limit defaults to 10",
			inputLastID:   &lastID,
			inputLimit:    -5,
			expectedLimit: 10,
		},
		{
			testName:      "success - nil lastID is accepted",
			inputLastID:   nil,
			inputLimit:    25,
			expectedLimit: 25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			pagination := domain.NewPagination(tc.inputLastID, tc.inputLimit)

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
		})
	}
}

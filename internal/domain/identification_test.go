package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewIdentification(t *testing.T) {
	testCases := []struct {
		testName      string
		idType        domain.IdentificationType
		number        string
		expectedError error
	}{
		{
			testName:      "success - create identification",
			idType:        domain.IdType_CC,
			number:        "123456789",
			expectedError: nil,
		},
		{
			testName:      "fail - invalid type",
			idType:        domain.IdentificationType("invalid"),
			number:        "123456789",
			expectedError: domain.ErrInvalidIdentificationType,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			id, err := domain.NewIdentification(tc.idType, tc.number)
			if err != tc.expectedError {
				t.Fatalf("expected error: %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				dto := id.ToDTO()
				if dto.Type != tc.idType || dto.Number != tc.number {
					t.Errorf("expected DTO: %+v", dto)
				}
			}
		})
	}
}

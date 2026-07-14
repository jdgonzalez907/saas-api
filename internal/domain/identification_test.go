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
				if id.Type() != tc.idType {
					t.Errorf("expected Type(): %s, got: %s", tc.idType, id.Type())
				}
				if id.Number() != tc.number {
					t.Errorf("expected Number(): %s, got: %s", tc.number, id.Number())
				}
			}
		})
	}
}

func TestIdentification_Equals(t *testing.T) {
	idBase, _ := domain.NewIdentification(domain.IdType_CC, "123456789")
	idSame, _ := domain.NewIdentification(domain.IdType_CC, "123456789")
	idDiffType, _ := domain.NewIdentification(domain.IdType_PASSPORT, "123456789")
	idDiffNum, _ := domain.NewIdentification(domain.IdType_CC, "987654321")

	testCases := []struct {
		testName string
		id1      domain.Identification
		id2      domain.Identification
		expected bool
	}{
		{
			testName: "success - identical identifications",
			id1:      idBase,
			id2:      idSame,
			expected: true,
		},
		{
			testName: "fail - different type",
			id1:      idBase,
			id2:      idDiffType,
			expected: false,
		},
		{
			testName: "fail - different number",
			id1:      idBase,
			id2:      idDiffNum,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			result := tc.id1.Equals(tc.id2)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

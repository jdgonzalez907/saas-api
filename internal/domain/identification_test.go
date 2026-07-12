package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewIdentification(t *testing.T) {
	testCases := []struct {
		testName       string
		input          domain.Identification
		expectedError  error
		expectedOutput domain.Identification
	}{
		{
			testName: "create identification",
			input: domain.Identification{
				Type:   domain.IdType_CC,
				Number: "123456789",
			},
			expectedError: nil,
			expectedOutput: domain.Identification{
				Type:   domain.IdType_CC,
				Number: "123456789",
			},
		},
		{
			testName: "fail to create identification with invalid type",
			input: domain.Identification{
				Type:   domain.IdentificationType("invalid"),
				Number: "123456789",
			},
			expectedError:  domain.ErrInvalidIdentificationType,
			expectedOutput: domain.Identification{},
		},
	}

	for _, testCase := range testCases {
		identification, err := domain.NewIdentification(testCase.input.Type, testCase.input.Number)
		if err != testCase.expectedError {
			t.Errorf("expected error: %v, got %v", testCase.expectedError, err)
		}
		if identification != testCase.expectedOutput {
			t.Errorf("expected identification: %v, got %v", testCase.expectedOutput, identification)
		}
	}
}

func TestIdentificationStrings(t *testing.T) {
	idType := domain.IdType_CC
	if idType.String() != "CC" {
		t.Errorf("expected %q, got %q", "CC", idType.String())
	}

	id := domain.Identification{
		Type:   domain.IdType_CC,
		Number: "12345",
	}
	expected := "CC12345"
	if id.String() != expected {
		t.Errorf("expected %q, got %q", expected, id.String())
	}
}

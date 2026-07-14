package domain_test

import (
	"testing"

	"jdgonzalez907/users-api/internal/domain"
)

func TestNewPhone(t *testing.T) {
	testCases := []struct {
		testName            string
		countryCode         string
		number              string
		expectedCountryCode string
		expectedNumber      string
		expectedError       error
	}{
		{
			testName:            "success - create phone",
			countryCode:         "57",
			number:              "3112223344",
			expectedCountryCode: "57",
			expectedNumber:      "3112223344",
			expectedError:       nil,
		},
		{
			testName:            "fail - empty country code",
			countryCode:         "",
			number:              "3112223344",
			expectedError:       domain.ErrInvalidPhone,
		},
		{
			testName:            "fail - empty number",
			countryCode:         "57",
			number:              "",
			expectedError:       domain.ErrInvalidPhone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			phone, err := domain.NewPhone(tc.countryCode, tc.number)
			if err != tc.expectedError {
				t.Fatalf("expected error: %v, got %v", tc.expectedError, err)
			}

			if tc.expectedError == nil {
				dto := phone.ToDTO()
				if dto.CountryCode != tc.expectedCountryCode {
					t.Errorf("expected DTO CountryCode: %s, got: %s", tc.expectedCountryCode, dto.CountryCode)
				}
				if dto.Number != tc.expectedNumber {
					t.Errorf("expected DTO Number: %s, got: %s", tc.expectedNumber, dto.Number)
				}
				if phone.CountryCode() != tc.expectedCountryCode {
					t.Errorf("expected CountryCode(): %s, got: %s", tc.expectedCountryCode, phone.CountryCode())
				}
				if phone.Number() != tc.expectedNumber {
					t.Errorf("expected Number(): %s, got: %s", tc.expectedNumber, phone.Number())
				}
			}
		})
	}
}

func TestPhone_Equals(t *testing.T) {
	phone1, _ := domain.NewPhone("57", "3001234567")
	phone2, _ := domain.NewPhone("57", "3001234567")
	phone3, _ := domain.NewPhone("1", "3001234567")
	phone4, _ := domain.NewPhone("57", "3009999999")

	if !phone1.Equals(phone2) {
		t.Error("expected phone1 to equal phone2")
	}
	if phone1.Equals(phone3) {
		t.Error("expected phone1 not to equal phone3")
	}
	if phone1.Equals(phone4) {
		t.Error("expected phone1 not to equal phone4")
	}
}

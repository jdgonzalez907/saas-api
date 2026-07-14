package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewPagination(t *testing.T) {
	now := time.Now().UTC()
	validID := int64(1)
	invalidID := int64(0)
	negativeID := int64(-5)

	limit10 := int32(10)
	limit25 := int32(25)
	limit50 := int32(50)
	invalidLimit := int32(15)

	testCases := []struct {
		name            string
		lastPublishedAt *time.Time
		lastID          *int64
		limitVal        *int32
		wantErr         error
	}{
		{
			name:            "success - all nil, default limit",
			lastPublishedAt: nil,
			lastID:          nil,
			limitVal:        nil,
			wantErr:         nil,
		},
		{
			name:            "success - all valid parameters",
			lastPublishedAt: &now,
			lastID:          &validID,
			limitVal:        &limit25,
			wantErr:         nil,
		},
		{
			name:            "success - allowed limit 10",
			lastPublishedAt: &now,
			lastID:          &validID,
			limitVal:        &limit10,
			wantErr:         nil,
		},
		{
			name:            "success - allowed limit 50",
			lastPublishedAt: &now,
			lastID:          &validID,
			limitVal:        &limit50,
			wantErr:         nil,
		},
		{
			name:            "fail - invalid limit",
			lastPublishedAt: nil,
			lastID:          nil,
			limitVal:        &invalidLimit,
			wantErr:         ErrInvalidPaginationLimit,
		},
		{
			name:            "fail - missing lastID when lastPublishedAt present",
			lastPublishedAt: &now,
			lastID:          nil,
			limitVal:        nil,
			wantErr:         ErrInvalidPaginationCursor,
		},
		{
			name:            "fail - missing lastPublishedAt when lastID present",
			lastPublishedAt: nil,
			lastID:          &validID,
			limitVal:        nil,
			wantErr:         ErrInvalidPaginationCursor,
		},
		{
			name:            "fail - invalid lastID (zero)",
			lastPublishedAt: &now,
			lastID:          &invalidID,
			limitVal:        nil,
			wantErr:         ErrInvalidPaginationCursor,
		},
		{
			name:            "fail - invalid lastID (negative)",
			lastPublishedAt: &now,
			lastID:          &negativeID,
			limitVal:        nil,
			wantErr:         ErrInvalidPaginationCursor,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewPagination(tc.lastPublishedAt, tc.lastID, tc.limitVal)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error %v, got %v", tc.wantErr, err)
			}

			if tc.wantErr == nil {
				// Verify getters
				if got.LastPublishedAt() != tc.lastPublishedAt {
					t.Errorf("expected lastPublishedAt %v, got %v", tc.lastPublishedAt, got.LastPublishedAt())
				}
				if got.LastID() != tc.lastID {
					t.Errorf("expected lastID %v, got %v", tc.lastID, got.LastID())
				}

				expectedLimit := DefaultLimit
				if tc.limitVal != nil {
					expectedLimit = *tc.limitVal
				}
				if got.Limit() != expectedLimit {
					t.Errorf("expected limit %d, got %d", expectedLimit, got.Limit())
				}
			}
		})
	}
}

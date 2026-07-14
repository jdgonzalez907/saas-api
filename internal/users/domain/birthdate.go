package domain

import (
	"errors"
	"time"
)

var (
	ErrUserUnderage           = errors.New("user must be at least 18 years old")
	ErrInvalidBirthDateFormat = errors.New("invalid birth date format")
)

const (
	MinAge          = 18
	BirthDateFormat = "2006-01-02"
)

type BirthDate struct {
	value     time.Time
	formatted string
}

func NewBirthDate(value string) (BirthDate, error) {
	t, err := time.Parse(BirthDateFormat, value)
	if err != nil {
		return BirthDate{}, ErrInvalidBirthDateFormat
	}

	now := time.Now().UTC()
	age := now.Year() - t.Year()

	if now.Month() < t.Month() || (now.Month() == t.Month() && now.Day() < t.Day()) {
		age--
	}
	if age < MinAge {
		return BirthDate{}, ErrUserUnderage
	}

	return BirthDate{
		value:     t,
		formatted: t.Format(BirthDateFormat),
	}, nil
}

type BirthDateDTO string

func (b BirthDate) Equals(other BirthDate) bool {
	return b.value.Equal(other.value) && b.formatted == other.formatted
}

func (b BirthDate) ToDTO() BirthDateDTO {
	return BirthDateDTO(b.formatted)
}

func (b BirthDate) Formatted() string {
	return b.formatted
}

func (b BirthDate) Time() time.Time {
	return b.value
}

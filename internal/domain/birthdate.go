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
	MinAge = 18
)

type BirthDate struct {
	value     time.Time
	formatted string
}

func NewBirthDate(value string) (BirthDate, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return BirthDate{}, ErrInvalidBirthDateFormat
	}

	now := time.Now().UTC()
	age := now.Year() - t.Year()

	if now.YearDay() < t.YearDay() {
		age--
	}
	if age < MinAge {
		return BirthDate{}, ErrUserUnderage
	}

	return BirthDate{
		value:     t,
		formatted: t.Format("2006-01-02"),
	}, nil
}

type BirthDateDTO string

func (b BirthDate) ToDTO() BirthDateDTO {
	return BirthDateDTO(b.formatted)
}

func (b BirthDate) Formatted() string {
	return b.formatted
}

func (b BirthDate) Time() time.Time {
	return b.value
}

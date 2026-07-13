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
	value time.Time
}

func NewBirthDate(value string) (BirthDate, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return BirthDate{}, ErrInvalidBirthDateFormat
	}

	now := time.Now()
	age := now.Year() - t.Year()

	if now.YearDay() < t.YearDay() {
		age--
	}
	if age < MinAge {
		return BirthDate{}, ErrUserUnderage
	}

	return BirthDate{value: t}, nil
}

type BirthDateDTO string

func (b BirthDate) ToDTO() BirthDateDTO {
	return BirthDateDTO(b.value.Format("2006-01-02"))
}

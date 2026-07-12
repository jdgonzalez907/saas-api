package domain

import (
	"errors"
	"time"
)

var ErrInvalidBirthDate = errors.New("invalid birth date")

const (
	MinAge = 18
)

type BirthDate struct {
	value time.Time
}

func NewBirthDate(value time.Time) (BirthDate, error) {
	now := time.Now()
	age := now.Year() - value.Year()

	if now.YearDay() < value.YearDay() {
		age--
	}
	if age < MinAge {
		return BirthDate{}, ErrInvalidBirthDate
	}

	return BirthDate{value: value}, nil
}

type BirthDateDTO struct {
	Value time.Time `json:"value"`
}

func (b BirthDate) ToDTO() BirthDateDTO {
	return BirthDateDTO{
		Value: b.value,
	}
}

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
	Value time.Time `json:"value"`
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

	return BirthDate{Value: value}, nil
}

func (b BirthDate) String() string {
	return b.Value.Format("2006-01-02")
}

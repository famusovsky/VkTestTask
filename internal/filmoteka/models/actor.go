package models

import (
	"errors"
	"time"
)

// ActorOut - структура, представляющая отправляемого актёра.
type ActorOut struct {
	Id          int       `json:"id" db:"id"`                       // Id - id актёра.
	Name        string    `json:"name" db:"name"`                   // Name - имя актёра.
	Gender      string    `json:"gender" db:"gender"`               // Gender - пол актёра.
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"` // DateOfBirth - дата рождения актёра.
	Movies      []int     `json:"movies" db:"-"`                    // Movies - список id фильмов, в которых принимал участие актёр.
}

// ActorIn - структура, представляющая получаемого актёра.
type ActorIn struct {
	Name        string    `json:"name" db:"name"`                   // Name - имя актёра.
	Gender      string    `json:"gender" db:"gender"`               // Gender - пол актёра.
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"` // DateOfBirth - дата рождения актёра.
}

// Check - проверка корректности данных актёра.
//
// Возвращает: ошибку.
func (a *ActorIn) Check() error {
	errs := make([]error, 0, 3)
	if a.Name == "" {
		errs = append(errs, errors.New("name must not be null"))
	}
	if a.Gender == "" {
		errs = append(errs, errors.New("gender must not be null"))
	}
	if a.DateOfBirth.IsZero() {
		errs = append(errs, errors.New("date of birth must not be null"))
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}

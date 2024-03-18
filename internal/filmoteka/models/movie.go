package models

import (
	"errors"
	"time"
)

// MovieOut - структура, представляющая отправляемый фильм.
type MovieOut struct {
	Id          int       `json:"id" db:"id"`                     // Id - id фильма.
	Name        string    `json:"name" db:"name"`                 // Name - название фильма.
	Description string    `json:"description" db:"description"`   // Description - описание фильма.
	ReleaseDate time.Time `json:"release_date" db:"release_date"` // ReleaseDate - дата выпуска фильма.
	Rating      int       `json:"rating" db:"rating"`             // Rating - рэйтинг фильма.
	Actors      []int     `json:"actors" db:"-"`                  // Actors - список id актёров, принимавших участие в фильме.
}

// MovieIn - структура, представляющая получаемый фильм.
type MovieIn struct {
	Name        string    `json:"name" db:"name"`                 // Name - название фильма.
	Description string    `json:"description" db:"description"`   // Description - описание фильма.
	ReleaseDate time.Time `json:"release_date" db:"release_date"` // ReleaseDate - дата выпуска фильма.
	Rating      *int      `json:"rating" db:"rating"`             // Rating - рэйтинг фильма.
	Actors      []int     `json:"actors" db:"-"`                  // Actors - список id актёров, принимавших участие в фильме.
}

// Check - проверка корректности данных фильма.
//
// Возвращает: ошибку.
func (m *MovieIn) Check() error {
	errs := make([]error, 0, 3)
	if m.Name == "" {
		errs = append(errs, errors.New("name must not be null"))
	}
	if len(m.Name) > 150 {
		errs = append(errs, errors.New("movie name must be less than 150 chars"))
	}
	if m.Description == "" {
		errs = append(errs, errors.New("description must not be null"))
	}
	if len(m.Description) > 1000 {
		errs = append(errs, errors.New("movie description must be less than 1000 chars"))
	}
	if m.ReleaseDate.IsZero() {
		errs = append(errs, errors.New("date of release must not be null"))
	}
	if m.Rating == nil {
		errs = append(errs, errors.New("rating of release must not be null"))
	} else if *m.Rating < 0 || *m.Rating > 10 {
		errs = append(errs, errors.New("rating of release must in range 0 - 10"))
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}

package models

import "time"

// Actor - структура, представляющая актёра.
type Actor struct {
	Id          int       `json:"id" db:"id"`                       // Id - id актёра.
	Name        string    `json:"name" db:"name"`                   // Name - имя актёра.
	Gender      string    `json:"gender" db:"gender"`               // Gender - пол актёра.
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"` // DateOfBirth - дата рождения актёра.
	Movies      []int     `json:"movies" db:"-"`                    // Movies - список id фильмов, в которых принимал участие актёр.
}

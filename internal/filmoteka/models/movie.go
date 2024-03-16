package models

import "time"

// Movie - структура, представляющая фильм.
type Movie struct {
	Id          int       `json:"id" db:"id"`                     // Id - id фильма.
	Name        string    `json:"name" db:"name"`                 // Name - название фильма.
	Description string    `json:"description" db:"description"`   // Description - описание фильма.
	ReleaseDate time.Time `json:"release_date" db:"release_date"` // ReleaseDate - дата выпуска фильма.
	Rating      int       `json:"rating" db:"rating"`             // Rating - рэйтинг фильма.
	Actors      []int     `json:"actors" db:"-"`                  // Actors - список id актёров, принимавших участие в фильме.
}

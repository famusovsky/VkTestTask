package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/famusovsky/VkTestTask/internal/filmoteka/models"
	"github.com/jmoiron/sqlx"
)

// TODO
type DbHandler interface {
	AddActor(name, gender string, dateOfBirth time.Time) (int, error)
	UpdateActor(id int, name, gender string, dateOfBirth time.Time) error
	DeleteActor(id int) error
	GetActor(id int) (models.Actor, error)
	GetActors() ([]models.Actor, error)

	AddMovie(name, description string, releaseDate time.Time, rating int, actors []int) (int, error)
	UpdateMovie(id int, name, description string, releaseDate time.Time, rating int, actors []int) error
	DeleteMovie(id int) error
	GetMovie(id int) (models.Movie, error)
	GetMovies(sortType int) ([]models.Movie, error)
	GetMoviesByActor(name string) ([]models.Movie, error)
	GetMoviesByName(name string) ([]models.Movie, error)

	AddUser(name, password string, isAdmin bool) (int, error)
	CheckUserRole(name, password string) (bool, error)
}

func GetHandler(db *sql.DB, overrideTables bool) (DbHandler, error) {
	if overrideTables {
		err := overrideDB(db)
		if err != nil {
			return dbProcessor{}, errors.Join(errors.New("error while getting db handler"), err)
		}
	}
	return dbProcessor{db: sqlx.NewDb(db, "postgres")}, nil
}

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

	AddMoovie(name, description string, releaseDate time.Time, rating int, actors []int) (int, error)
	UpdateMoovie(id int, name, description string, releaseDate time.Time, rating int, actors []int) error
	DeleteMoovie(id int) error
	GetMoovie(id int) (models.Moovie, error)
	GetMoovies(sortType int) ([]models.Moovie, error)
	GetMooviesByActor(name string) ([]models.Moovie, error)
	GetMooviesByName(name string) ([]models.Moovie, error)

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

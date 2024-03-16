package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/famusovsky/VkTestTask/internal/filmoteka/models"
	"github.com/jmoiron/sqlx"
)

// TODO refactor

// overrideDB - функция, перезаписывающая таблицы фильмотеки в БД.
func overrideDB(db *sql.DB) error {
	wrapErr := errors.New("error while overriding db")
	err := dropTables(db)
	if err != nil {
		return errors.Join(wrapErr, err)
	}
	err = createTables(db)
	if err != nil {
		return errors.Join(wrapErr, err)
	}
	return nil
}

// dropTables - функция, удаляющая таблицы фильмотеки в БД.
func dropTables(db *sql.DB) error {
	q := strings.Join([]string{dropMovieActors, dropActors, dropMovies, dropUsers}, " ")

	_, err := db.Exec(q)
	if err != nil {
		return errors.Join(fmt.Errorf("error while dropping tables: %s", err))
	}

	return nil
}

// createTables - функция, добавляющая таблицы фильмотеки в БД.
func createTables(db *sql.DB) error {
	q := strings.Join([]string{createActors, createMovies, createUsers, createActorMovieRelations}, " ")

	_, err := db.Exec(q)
	if err != nil {
		return errors.Join(fmt.Errorf("error while creating tables: %s", err))
	}

	return nil
}

// dbProcessor - структура, представляющая обработчик БД.
type dbProcessor struct {
	db *sqlx.DB
}

var (
	// Ошибка создания SQL транзакции.
	errBeginTx = errors.New("error while starting transaction")
	// Ошибка сохранения SQL транзакции.
	errCommitTx = errors.New("error while committing transaction")
)

func (d dbProcessor) addSmthWithId(query string, wrap string, args ...any) (int, error) {
	wrapErr := errors.New(wrap)
	tx, err := d.db.Begin()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.Join(wrapErr, errCommitTx, err)
	}

	return id, nil
}

// AddActor implements DbHandler.
func (d dbProcessor) AddActor(name string, gender string, dateOfBirth time.Time) (int, error) {
	return d.addSmthWithId(addActor,
		"error while inserting actor",
		name, gender, dateOfBirth)
}

// AddUser implements DbHandler.
func (d dbProcessor) AddUser(name string, password string, isAdmin bool) (int, error) {
	return d.addSmthWithId(addUser,
		"error while inserting user",
		name, password, isAdmin)
}

// AddMovie implements DbHandler.
func (d dbProcessor) AddMovie(name string, description string, releaseDate time.Time, rating int, actors []int) (int, error) {
	wrapErr := errors.New("error while inserting actor")
	tx, err := d.db.Begin()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRow(addMovie, name, description, releaseDate, rating).Scan(&id)
	if err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	for _, actorId := range actors {
		err := d.addActorToMovie(tx, actorId, id)
		if err != nil {
			return 0, errors.Join(wrapErr, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.Join(wrapErr, errCommitTx, err)
	}

	return id, nil
}

func (d dbProcessor) addActorToMovie(tx *sql.Tx, actorId, movieId int) error {
	_, err := tx.Exec(addActorToMovie, movieId, actorId)
	if err != nil {
		return errors.Join(fmt.Errorf("error while adding actor %d to movie %d", actorId, movieId), err)
	}

	return nil
}

// CheckUserRole implements DbHandler.
func (d dbProcessor) CheckUserRole(name string, password string) (bool, error) {
	var isAdmin bool
	err := d.db.Get(&isAdmin, checkUserRole, name, password)
	if err != nil {
		return false, errors.Join(errors.New("error while user's role"), err)
	}

	return isAdmin, nil
}

func (d dbProcessor) deleteSmth(query, errTxt string, args ...any) error {
	wrapErr := errors.New(errTxt)
	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// DeleteActor implements DbHandler.
func (d dbProcessor) DeleteActor(id int) error {
	return d.deleteSmth(removeActor, fmt.Sprintf("error while deleting actor %d", id))
}

// DeleteMovie implements DbHandler.
func (d dbProcessor) DeleteMovie(id int) error {
	return d.deleteSmth(removeMovie, fmt.Sprintf("error while deleting movie %d", id))
}

// GetActor implements DbHandler.
func (d dbProcessor) GetActor(id int) (models.Actor, error) {
	var actor models.Actor
	err := d.db.Get(&actor, getActor, id)
	if err != nil {
		return models.Actor{}, errors.Join(errors.New("error while getting actor"), err)
	}

	return actor, nil
}

// GetActors implements DbHandler.
func (d dbProcessor) GetActors() ([]models.Actor, error) {
	var actors []models.Actor
	err := d.db.Select(&actors, getActors)
	if err != nil {
		return nil, errors.Join(errors.New("error while getting actors"), err)
	}

	return actors, nil
}

// GetMovie implements DbHandler.
func (d dbProcessor) GetMovie(id int) (models.Movie, error) {
	wrapErr := fmt.Errorf("error while getting movie %d", id)
	var movie models.Movie
	err := d.db.Select(&movie, getMovie)
	if err != nil {
		return models.Movie{}, errors.Join(wrapErr, err)
	}

	err = d.db.Select(&movie.Actors, getMovieActors, id)
	if err != nil {
		return models.Movie{}, errors.Join(wrapErr, err)
	}

	return movie, nil
}

func (d dbProcessor) fillMovies(movies []models.Movie) error {
	for i := range len(movies) {
		err := d.db.Select(&movies[i].Actors, getMovieActors, movies[i].Id)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetMovies implements DbHandler.
func (d dbProcessor) GetMovies(sortType int) ([]models.Movie, error) {
	wrapErr := errors.New("error while getting movies")
	var movies []models.Movie
	var err error
	switch sortType {
	case models.SortByRating:
		err = d.db.Select(&movies, getMoviesSortByRating)
	case models.SortByName:
		err = d.db.Select(&movies, getMoviesSortByName)
	case models.SortByReleaseDate:
		err = d.db.Select(&movies, getMoviesSortByReleaseDate)
	}
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	err = d.fillMovies(movies)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	return movies, nil
}

// GetMoviesByActor implements DbHandler.
func (d dbProcessor) GetMoviesByActor(name string) ([]models.Movie, error) {
	wrapErr := errors.New("error while getting movies by actor")
	var movies []models.Movie
	err := d.db.Select(&movies, getMoviesByActor, name)

	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	err = d.fillMovies(movies)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	return movies, nil
}

// GetMoviesByName implements DbHandler.
func (d dbProcessor) GetMoviesByName(name string) ([]models.Movie, error) {
	wrapErr := errors.New("error while getting movies by name")
	var movies []models.Movie
	err := d.db.Select(&movies, getMoviesByName, name)

	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	err = d.fillMovies(movies)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	return movies, nil
}

// UpdateActor implements DbHandler.
func (d dbProcessor) UpdateActor(id int, name string, gender string, dateOfBirth time.Time) error {
	wrapErr := fmt.Errorf("error while updating actor %d", id)
	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(updateActor, id, name, gender, dateOfBirth)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return err
}

// UpdateMovie implements DbHandler.
func (d dbProcessor) UpdateMovie(id int, name string, description string, releaseDate time.Time, rating int, actors []int) error {
	wrapErr := fmt.Errorf("error while updating movie %d", id)
	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(updateMovie, id, name, description, releaseDate, rating)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	_, err = tx.Exec(removeMovieFromActors, id)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	for _, actorId := range actors {
		err := d.addActorToMovie(tx, actorId, id)
		if err != nil {
			return errors.Join(wrapErr, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return err
}

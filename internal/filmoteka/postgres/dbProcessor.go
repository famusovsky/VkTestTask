package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/famusovsky/VkTestTask/internal/filmoteka/models"
	"github.com/jmoiron/sqlx"
)

// overrideDB - функция, перезаписывающая таблицы фильмотеки в БД.
func overrideDB(db *sql.DB) error {
	wrapErr := errors.New("error while overriding db")
	if err := dropTables(db); err != nil {
		return errors.Join(wrapErr, err)
	}
	if err := createTables(db); err != nil {
		return errors.Join(wrapErr, err)
	}
	return nil
}

// dropTables - функция, удаляющая таблицы фильмотеки в БД.
func dropTables(db *sql.DB) error {
	q := strings.Join([]string{dropMovieActors, dropActors, dropMovies, dropUsers}, " ")
	if _, err := db.Exec(q); err != nil {
		return errors.Join(fmt.Errorf("error while dropping tables: %s", err))
	}
	return nil
}

// createTables - функция, добавляющая таблицы фильмотеки в БД.
func createTables(db *sql.DB) error {
	q := strings.Join([]string{createActors, createMovies, createUsers, createActorMovieRelations}, " ")
	if _, err := db.Exec(q); err != nil {
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

// AddActor - добавление актёра в БД.
func (d dbProcessor) AddActor(a models.ActorIn) (int, error) {
	return d.addSmthWithId(addActor, "error while inserting actor", a.Name, a.Gender, a.DateOfBirth)
}

// AddUser - добавление пользователя в БД.
func (d dbProcessor) AddUser(u models.User) (int, error) {
	return d.addSmthWithId(addUser, "error while inserting user", u.Nickname, u.Password, u.IsAdmin)
}

// AddMovie - добавление фильма в БД.
func (d dbProcessor) AddMovie(m models.MovieIn) (int, error) {
	wrapErr := errors.New("error while inserting movie")
	tx, err := d.db.Beginx()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int
	if err = tx.QueryRow(addMovie, m.Name, m.Description, m.ReleaseDate, *m.Rating).Scan(&id); err != nil {
		return 0, errors.Join(wrapErr, err)
	}
	for _, aId := range m.Actors {
		err = errors.Join(d.addActorToMovie(tx, aId, id))
		if err != nil {
			return 0, errors.Join(wrapErr, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, errors.Join(wrapErr, errCommitTx, err)
	}

	return id, err
}

// CheckUserRole - проверка роли пользователя.
func (d dbProcessor) CheckUserRole(name string, password string) (bool, error) {
	var isAdmin bool
	if err := d.db.Get(&isAdmin, checkUserRole, name, password); err != nil {
		return false, errors.Join(errors.New("error while checking user's role"), err)
	}
	return isAdmin, nil
}

// DeleteActor - удаление актёра из БД.
func (d dbProcessor) DeleteActor(id int) error {
	return d.deleteSmth(removeActor, fmt.Sprintf("error while deleting actor %d", id), id)
}

// DeleteMovie - удаление фильма из БД.
func (d dbProcessor) DeleteMovie(id int) error {
	return d.deleteSmth(removeMovie, fmt.Sprintf("error while deleting movie %d", id), id)
}

// GetActor - получение актёра из БД.
func (d dbProcessor) GetActor(id int) (models.ActorOut, error) {
	wrapErr := fmt.Errorf("error while getting actor %d", id)
	var actor models.ActorOut
	if err := d.db.Get(&actor, getActor, id); err != nil {
		return models.ActorOut{}, errors.Join(wrapErr, err)
	}
	if err := d.db.Select(&actor.Movies, getActorMovies, id); err != nil {
		return actor, errors.Join(wrapErr, errors.New("error while getting actors's movies"), err)
	}
	return actor, nil
}

// GetActors - получение актёров из БД.
func (d dbProcessor) GetActors() ([]models.ActorOut, error) {
	var actors []models.ActorOut
	if err := d.db.Select(&actors, getActors); err != nil {
		return nil, errors.Join(errors.New("error while getting actors"), err)
	}
	for i := range actors {
		if err := d.db.Select(&actors[i].Movies, getActorMovies, actors[i].Id); err != nil {
			return nil, errors.Join(errors.New("error while getting actors' movies"), err)
		}
	}
	return actors, nil
}

// GetMovie - получение фильма из БД.
func (d dbProcessor) GetMovie(id int) (models.MovieOut, error) {
	wrapErr := fmt.Errorf("error while getting movie %d", id)
	var movie models.MovieOut
	if err := d.db.Get(&movie, getMovie, id); err != nil {
		return models.MovieOut{}, errors.Join(wrapErr, err)
	}
	if err := d.db.Select(&movie.Actors, getMovieActors, id); err != nil {
		return movie, errors.Join(wrapErr, errors.New("error while getting movie's actors"), err)
	}
	return movie, nil
}

// GetMovies - получение фильмов из БД.
func (d dbProcessor) GetMovies(sortType int) ([]models.MovieOut, error) {
	wrapErr := errors.New("error while getting movies")
	var movies []models.MovieOut
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

	if err = d.fillMovies(movies); err != nil {
		return nil, errors.Join(wrapErr, errors.New("error while getting movie's actors"), err)
	}
	return movies, nil
}

// GetMoviesByActor - получение фильмов, в которых играл актёр, из БД.
func (d dbProcessor) GetMoviesByActor(name string) ([]models.MovieOut, error) {
	wrapErr := errors.New("error while getting movies by actor")
	var movies []models.MovieOut
	if err := d.db.Select(&movies, getMoviesByActor, name); err != nil {
		return nil, errors.Join(wrapErr, err)
	}
	if err := d.fillMovies(movies); err != nil {
		return nil, errors.Join(wrapErr, errors.New("error while getting movie's actors"), err)
	}
	return movies, nil
}

// GetMoviesByName - получение фильмов по фрагменту названия из БД.
func (d dbProcessor) GetMoviesByName(name string) ([]models.MovieOut, error) {
	wrapErr := errors.New("error while getting movies by name")
	var movies []models.MovieOut
	if err := d.db.Select(&movies, getMoviesByName, name); err != nil {
		return nil, errors.Join(wrapErr, err)
	}
	if err := d.fillMovies(movies); err != nil {
		return nil, errors.Join(wrapErr, err)
	}
	return movies, nil
}

// UpdateActor - обновление актёра в БД.
func (d dbProcessor) UpdateActor(id int, a models.ActorIn) error {
	wrapErr := fmt.Errorf("error while updating actor %d", id)
	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if a.Name != "" {
		if _, err = tx.Exec(updateActorName, id, a.Name); err != nil {
			return errors.Join(wrapErr, err)
		}
	}
	if a.Gender != "" {
		if _, err = tx.Exec(updateActorGender, id, a.Gender); err != nil {
			return errors.Join(wrapErr, err)
		}
	}
	if !a.DateOfBirth.IsZero() {
		if _, err = tx.Exec(updateActorDateOfBirth, id, a.DateOfBirth); err != nil {
			return errors.Join(wrapErr, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}
	return nil
}

// UpdateMovie - обновление фильма в БД.
func (d dbProcessor) UpdateMovie(id int, m models.MovieIn) error {
	wrapErr := fmt.Errorf("error while updating movie %d", id)
	tx, err := d.db.Beginx()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if m.Name != "" {
		if _, err = tx.Exec(updateMovieName, id, m.Name); err != nil {
			return errors.Join(wrapErr, err)
		}
	}
	if m.Description != "" {
		if _, err = tx.Exec(updateMovieDescription, id, m.Description); err != nil {
			return errors.Join(wrapErr, err)
		}
	}
	if !m.ReleaseDate.IsZero() {
		if _, err = tx.Exec(updateMovieReleaseDate, id, m.ReleaseDate); err != nil {
			return errors.Join(wrapErr, err)
		}
	}
	if m.Rating != nil {
		if _, err = tx.Exec(updateMovieRating, id, *m.Rating); err != nil {
			return errors.Join(wrapErr, err)
		}
	}

	if m.Actors != nil {
		if _, err = tx.Exec(removeMovieFromActors, id); err != nil {
			return errors.Join(wrapErr, err)
		}
		for _, aId := range m.Actors {
			err = errors.Join(d.addActorToMovie(tx, aId, id))
			if err != nil {
				return errors.Join(wrapErr, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// addActorToMovie - добавление актёра в фильм.
func (d dbProcessor) addActorToMovie(tx *sqlx.Tx, actorId, movieId int) error {
	_, err := tx.Exec(addActorToMovie, movieId, actorId)
	if err != nil {
		return errors.Join(fmt.Errorf("error while adding actor %d to movie %d", actorId, movieId), err)
	}

	return nil
}

// addSmthWithId - добавление чего-либо в БД с возвращением id.
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

// deleteSmth - удаление чего-либо из БД.
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

// fillMovies - заполнение фильмов актёрами.
func (d dbProcessor) fillMovies(movies []models.MovieOut) error {
	for i := range len(movies) {
		err := d.db.Select(&movies[i].Actors, getMovieActors, movies[i].Id)
		if err != nil {
			return err
		}
	}

	return nil
}

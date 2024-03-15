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
	q := strings.Join([]string{dropMoovieActors, dropActors, dropMoovies, dropUsers}, " ")

	_, err := db.Exec(q)
	if err != nil {
		return errors.Join(fmt.Errorf("error while dropping tables: %s", err))
	}

	return nil
}

// createTables - функция, добавляющая таблицы фильмотеки в БД.
func createTables(db *sql.DB) error {
	q := strings.Join([]string{createActors, createMoovies, createUsers, createActorMoovieRelations}, " ")

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

// AddMoovie implements DbHandler.
func (d dbProcessor) AddMoovie(name string, description string, releaseDate time.Time, rating int, actors []int) (int, error) {
	wrapErr := errors.New("error while inserting actor")
	tx, err := d.db.Begin()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRow(addMoovie, name, description, releaseDate, rating).Scan(&id)
	if err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	for _, actorId := range actors {
		err := d.addActorToMoovie(tx, actorId, id)
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

func (d dbProcessor) addActorToMoovie(tx *sql.Tx, actorId, moovieId int) error {
	_, err := tx.Exec(addActorToMoovie, moovieId, actorId)
	if err != nil {
		return errors.Join(fmt.Errorf("error while adding actor %d to moovie %d", actorId, moovieId), err)
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

// DeleteMoovie implements DbHandler.
func (d dbProcessor) DeleteMoovie(id int) error {
	return d.deleteSmth(removeMoovie, fmt.Sprintf("error while deleting moovie %d", id))
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

// GetMoovie implements DbHandler.
func (d dbProcessor) GetMoovie(id int) (models.Moovie, error) {
	wrapErr := fmt.Errorf("error while getting moovie %d", id)
	var moovie models.Moovie
	err := d.db.Select(&moovie, getMoovie)
	if err != nil {
		return models.Moovie{}, errors.Join(wrapErr, err)
	}

	err = d.db.Select(&moovie.Actors, getMoovieActors, id)
	if err != nil {
		return models.Moovie{}, errors.Join(wrapErr, err)
	}

	return moovie, nil
}

func (d dbProcessor) fillMoovies(moovies []models.Moovie) error {
	for i := range len(moovies) {
		err := d.db.Select(&moovies[i].Actors, getMoovieActors, moovies[i].Id)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetMoovies implements DbHandler.
func (d dbProcessor) GetMoovies(sortType int) ([]models.Moovie, error) {
	wrapErr := errors.New("error while getting moovies")
	var moovies []models.Moovie
	var err error
	switch sortType {
	case models.SortByRating:
		err = d.db.Select(&moovies, getMooviesSortByRating)
	case models.SortByName:
		err = d.db.Select(&moovies, getMooviesSortByName)
	case models.SortByReleaseDate:
		err = d.db.Select(&moovies, getMooviesSortByReleaseDate)
	}
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	err = d.fillMoovies(moovies)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	return moovies, nil
}

// GetMooviesByActor implements DbHandler.
func (d dbProcessor) GetMooviesByActor(name string) ([]models.Moovie, error) {
	wrapErr := errors.New("error while getting moovies by actor")
	var moovies []models.Moovie
	err := d.db.Select(&moovies, getMooviesByActor, name)

	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	err = d.fillMoovies(moovies)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	return moovies, nil
}

// GetMooviesByName implements DbHandler.
func (d dbProcessor) GetMooviesByName(name string) ([]models.Moovie, error) {
	wrapErr := errors.New("error while getting moovies by name")
	var moovies []models.Moovie
	err := d.db.Select(&moovies, getMooviesByName, name)

	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	err = d.fillMoovies(moovies)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	return moovies, nil
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

// UpdateMoovie implements DbHandler.
func (d dbProcessor) UpdateMoovie(id int, name string, description string, releaseDate time.Time, rating int, actors []int) error {
	wrapErr := fmt.Errorf("error while updating moovie %d", id)
	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(updateMoovie, id, name, description, releaseDate, rating)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	_, err = tx.Exec(removeMoovieFromActors, id)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	for _, actorId := range actors {
		err := d.addActorToMoovie(tx, actorId, id)
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

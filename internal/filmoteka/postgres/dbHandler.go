package postgres

import (
	"database/sql"
	"errors"

	"github.com/famusovsky/VkTestTask/internal/filmoteka/models"
	"github.com/jmoiron/sqlx"
)

// DbHandler - обработчик базы данных фильмотеки.
type DbHandler interface {
	// AddActor - добавляет актера в базу данных.
	//
	// Принимает: актёр.
	//
	// Возвращает: id добавленного актёра и ошибку.
	AddActor(a models.ActorIn) (int, error)

	// UpdateActor - обновляет актёра в базе данных.
	//
	// Принимает: id актёра и обновлённые данные актёра.
	//
	// Возвращает: ошибку.
	UpdateActor(id int, a models.ActorIn) error

	// DeleteActor - удаляет актёра из базы данных.
	//
	// Принимает: id актёра.
	//
	// Возвращает: ошибку.
	DeleteActor(id int) error

	// GetActor - получает актёра из базы данных.
	//
	// Принимает: id актёра.
	//
	// Возвращает: актёра и ошибку.
	GetActor(id int) (models.ActorOut, error)

	// GetActors - получает всех актёров из базы данных.
	//
	// Возвращает: всех актёров и ошибку.
	GetActors() ([]models.ActorOut, error)

	// AddMovie - добавляет фильм в базу данных.
	//
	// Принимает: фильм.
	//
	// Возвращает: id добавленного фильма и ошибку.
	AddMovie(m models.MovieIn) (int, error)

	// UpdateMovie - обновляет фильм в базе данных.
	//
	// Принимает: id фильма и обновлённые данные фильма.
	//
	// Возвращает: ошибку.
	UpdateMovie(id int, m models.MovieIn) error

	// DeleteMovie - удаляет фильм из базы данных.
	//
	// Принимает: id фильма.
	//
	// Возвращает: ошибку.
	DeleteMovie(id int) error

	// GetMovie - получает фильм из базы данных.
	//
	// Принимает: id фильма.
	//
	// Возвращает: фильм и ошибку.
	GetMovie(id int) (models.MovieOut, error)

	// GetMovies - получает все фильмы из базы данных.
	//
	// Принимает: тип сортировки.
	//
	// Возвращает: все фильмы и ошибку.
	GetMovies(sortType int) ([]models.MovieOut, error)

	// GetMoviesByActor - получает все фильмы с участием актёра из базы данных.
	//
	// Принимает: имя актёра.
	//
	// Возвращает: все фильмы и ошибку.
	GetMoviesByActor(name string) ([]models.MovieOut, error)

	// GetMoviesByName - получает все фильмы с именем из базы данных.
	//
	// Принимает: имя фильма.
	//
	// Возвращает: все фильмы и ошибку.
	GetMoviesByName(name string) ([]models.MovieOut, error)

	// AddUser - добавляет пользователя в базу данных.
	//
	// Принимает: пользователя.
	//
	// Возвращает: id добавленного пользователя и ошибку.
	AddUser(u models.User) (int, error)

	// UpdateUser - обновляет пользователя в базе данных.
	//
	// Принимает: id пользователя и обновлённые данные пользователя.
	//
	// Возвращает: ошибку.
	CheckUserRole(name, password string) (bool, error)
}

// GetHandler - возвращает обработчик базы данных фильмотеки.
//
// Принимает: подключение к базе данных и флаг пересоздания таблиц.
//
// Возвращает: обработчик базы данных и ошибку.
func GetHandler(db *sql.DB, overrideTables bool) (DbHandler, error) {
	if overrideTables {
		err := overrideDB(db)
		if err != nil {
			return dbProcessor{}, errors.Join(errors.New("error while getting db handler"), err)
		}
	}
	return dbProcessor{db: sqlx.NewDb(db, "postgres")}, nil
}

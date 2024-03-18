package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/famusovsky/VkTestTask/internal/filmoteka/models"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestOverrideDB(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	mock.MatchExpectationsInOrder(false)

	mock.ExpectExec(dropMovieActors).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(dropActors).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(dropMovies).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(dropUsers).WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS actors").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS movies").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS movie_actors").WillReturnResult(sqlmock.NewResult(0, 0))

	err = overrideDB(mockDB)
	assert.NoError(t, err)
}

func TestDropTables(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		mock.ExpectExec(dropMovieActors).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(dropActors).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(dropMovies).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(dropUsers).WillReturnResult(sqlmock.NewResult(0, 0))

		err = dropTables(mockDB)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		errTxt := "drop error"
		mock.ExpectExec(dropMovieActors).WillReturnError(errors.New(errTxt))
		mock.ExpectExec(dropActors).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(dropMovies).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(dropUsers).WillReturnResult(sqlmock.NewResult(0, 0))

		err = dropTables(mockDB)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while dropping tables:")
	})
}

func TestCreateTables(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS actors").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS movies").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS movie_actors").WillReturnResult(sqlmock.NewResult(0, 0))

		err = createTables(mockDB)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		errTxt := "create error"
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS actors").WillReturnError(errors.New(errTxt))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS movies").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("CREATE TABLE IF NOT EXISTS movie_actors").WillReturnResult(sqlmock.NewResult(0, 0))

		err = createTables(mockDB)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while creating tables:")
	})
}

func TestAddSmthWithId(t *testing.T) {
	smth := true
	q := "INSERT INTO smth"
	wrap := "wrap"

	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		mock.ExpectBegin()
		mock.ExpectQuery(q).WithArgs(smth).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		id, err := processor.addSmthWithId(q, wrap, smth)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("begin error", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		mock.ExpectBegin().WillReturnError(errors.New("begin error"))
		_, err := processor.addSmthWithId(q, wrap, smth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errBeginTx.Error())
		assert.Contains(t, err.Error(), wrap)
	})

	t.Run("insert error", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		mock.ExpectBegin()
		mock.ExpectQuery(q).WithArgs(smth).WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		_, err := processor.addSmthWithId(q, wrap, smth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), wrap)
		assert.Contains(t, err.Error(), "insert error")
	})

	t.Run("commit error", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		mock.ExpectBegin()
		mock.ExpectQuery(q).WithArgs(smth).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		_, err := processor.addSmthWithId(q, wrap, smth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), wrap)
		assert.Contains(t, err.Error(), errCommitTx.Error())
		assert.Contains(t, err.Error(), "commit error")
	})
}

func TestAddActor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		actor := models.ActorIn{}
		id := 15
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO actors").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
		mock.ExpectCommit()

		id, err := processor.AddActor(actor)
		assert.NoError(t, err)
		assert.Equal(t, id, id)
	})
}

func TestAddUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		user := models.User{}
		id := 15
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO user").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
		mock.ExpectCommit()

		id, err := processor.AddUser(user)
		assert.NoError(t, err)
		assert.Equal(t, id, id)
	})
}

func TestAddMovie(t *testing.T) {
	var (
		m = models.MovieIn{
			Name:        "title",
			Description: "description",
			ReleaseDate: time.Time{},
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		q1   = "INSERT INTO movies"
		q2   = "INSERT INTO movie_actors"
		wrap = "error while inserting movie"
	)
	*m.Rating = 5

	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		mock.ExpectBegin()
		mock.ExpectQuery(q1).WithArgs(m.Name, m.Description, m.ReleaseDate, m.Rating).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		for _, a := range m.Actors {
			mock.ExpectExec(q2).WithArgs(1, a).WillReturnResult(sqlmock.NewResult(0, 0))
		}
		mock.ExpectCommit()

		id, err := processor.AddMovie(m)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("error while starting", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		errTxt := "begin error"
		mock.ExpectBegin().WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		_, err := processor.AddMovie(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), wrap)
		assert.Contains(t, err.Error(), errTxt)
	})

	t.Run("error while inserting movie", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		errTxt := "insert error"
		mock.ExpectBegin()
		mock.ExpectQuery(q1).WithArgs(m.Name, m.Description, m.ReleaseDate, m.Rating).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		_, err := processor.AddMovie(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), wrap)
		assert.Contains(t, err.Error(), errTxt)
	})

	t.Run("error while inserting actors", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		errTxt := "insert error"
		mock.ExpectBegin()
		mock.ExpectQuery(q1).WithArgs(m.Name, m.Description, m.ReleaseDate, m.Rating).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectExec(q2).WithArgs(1, m.Actors[0]).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(q2).WithArgs(1, m.Actors[1]).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		_, err := processor.AddMovie(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), wrap)
		assert.Contains(t, err.Error(), errTxt)
	})

	t.Run("error while committing", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		errTxt := "commit error"
		mock.ExpectBegin()
		mock.ExpectQuery(q1).WithArgs(m.Name, m.Description, m.ReleaseDate, m.Rating).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		for _, a := range m.Actors {
			mock.ExpectExec(q2).WithArgs(1, a).WillReturnResult(sqlmock.NewResult(0, 0))
		}
		mock.ExpectCommit().WillReturnError(errors.New(errTxt))

		_, err := processor.AddMovie(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), wrap)
		assert.Contains(t, err.Error(), errTxt)
	})
}

func TestGetActor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		id := 15
		actor := models.ActorOut{
			Id:          id,
			Name:        "name",
			DateOfBirth: time.Time{},
			Movies:      []int{1},
		}

		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "date_of_birth"}).AddRow(actor.Id, actor.Name, actor.DateOfBirth))
		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		a, err := processor.GetActor(id)
		assert.NoError(t, err)
		assert.Equal(t, actor, a)
	})

	t.Run("error", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		id := 15
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(errors.New(errTxt))

		_, err := processor.GetActor(id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting actor 15")
	})

	t.Run("error while getting movies", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		id := 15
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "date_of_birth"}).AddRow(id, "name", time.Time{}))
		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(errors.New(errTxt))

		_, err := processor.GetActor(id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting actor 15")
		assert.Contains(t, err.Error(), "error while getting actors's movies")
	})
}

func TestGetActors(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		actors := []models.ActorOut{
			{
				Id:          1,
				Name:        "name1",
				DateOfBirth: time.Time{},
				Movies:      []int{1},
			},
			{
				Id:          2,
				Name:        "name2",
				DateOfBirth: time.Time{},
				Movies:      []int{2},
			},
		}

		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "date_of_birth"}).AddRow(actors[0].Id, actors[0].Name, actors[0].DateOfBirth).AddRow(actors[1].Id, actors[1].Name, actors[1].DateOfBirth))
		mock.ExpectQuery("SELECT").WithArgs(actors[0].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT").WithArgs(actors[1].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		a, err := processor.GetActors()
		assert.NoError(t, err)
		assert.Equal(t, actors, a)
	})

	t.Run("error while getting actors", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WillReturnError(errors.New(errTxt))

		_, err := processor.GetActors()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting actors")
	})

	t.Run("error while getting movies", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "date_of_birth"}).AddRow(1, "name", time.Time{}))
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New(errTxt))

		_, err := processor.GetActors()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting actors' movies")
	})
}

func TestGetMovie(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		id := 15
		movie := models.MovieOut{
			Id:          id,
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Time{},
			Rating:      5,
			Actors:      []int{1},
		}

		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(movie.Id, movie.Name, movie.Description, movie.ReleaseDate, movie.Rating))
		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		m, err := processor.GetMovie(id)
		assert.NoError(t, err)
		assert.Equal(t, movie, m)
	})

	t.Run("error while getting movie", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		id := 15
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(errors.New(errTxt))

		_, err := processor.GetMovie(id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting movie")
	})

	t.Run("error while getting actors", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		id := 15
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(id, "name", "description", time.Time{}, 5))
		mock.ExpectQuery("SELECT").WithArgs(id).WillReturnError(errors.New(errTxt))

		_, err := processor.GetMovie(id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting movie's actors")
		assert.Contains(t, err.Error(), "error while getting movie")
	})
}

func TestGetMovies(t *testing.T) {
	t.Run("success sort by rating", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		movies := []models.MovieOut{
			{
				Id:          1,
				Name:        "name1",
				Description: "description1",
				ReleaseDate: time.Time{},
				Rating:      5,
				Actors:      []int{1},
			},
			{
				Id:          2,
				Name:        "name2",
				Description: "description2",
				ReleaseDate: time.Time{},
				Rating:      4,
				Actors:      []int{2},
			},
		}

		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(movies[0].Id, movies[0].Name, movies[0].Description, movies[0].ReleaseDate, movies[0].Rating).AddRow(movies[1].Id, movies[1].Name, movies[1].Description, movies[1].ReleaseDate, movies[1].Rating))
		mock.ExpectQuery("SELECT").WithArgs(movies[0].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT").WithArgs(movies[1].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		m, err := processor.GetMovies(models.SortByRating)
		assert.NoError(t, err)
		assert.Equal(t, movies, m)
	})

	t.Run("success sort by name", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		movies := []models.MovieOut{
			{
				Id:          1,
				Name:        "name1",
				Description: "description1",
				ReleaseDate: time.Time{},
				Rating:      5,
				Actors:      []int{1},
			},
			{
				Id:          2,
				Name:        "name2",
				Description: "description2",
				ReleaseDate: time.Time{},
				Rating:      4,
				Actors:      []int{2},
			},
		}

		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(movies[0].Id, movies[0].Name, movies[0].Description, movies[0].ReleaseDate, movies[0].Rating).AddRow(movies[1].Id, movies[1].Name, movies[1].Description, movies[1].ReleaseDate, movies[1].Rating))
		mock.ExpectQuery("SELECT").WithArgs(movies[0].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT").WithArgs(movies[1].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		m, err := processor.GetMovies(models.SortByName)
		assert.NoError(t, err)
		assert.Equal(t, movies, m)
	})

	t.Run("success sort by release date", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		movies := []models.MovieOut{
			{
				Id:          1,
				Name:        "name1",
				Description: "description1",
				ReleaseDate: time.Time{},
				Rating:      5,
				Actors:      []int{1},
			},
			{
				Id:          2,
				Name:        "name2",
				Description: "description2",
				ReleaseDate: time.Time{},
				Rating:      4,
				Actors:      []int{2},
			},
		}

		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(movies[0].Id, movies[0].Name, movies[0].Description, movies[0].ReleaseDate, movies[0].Rating).AddRow(movies[1].Id, movies[1].Name, movies[1].Description, movies[1].ReleaseDate, movies[1].Rating))
		mock.ExpectQuery("SELECT").WithArgs(movies[0].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT").WithArgs(movies[1].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		m, err := processor.GetMovies(models.SortByReleaseDate)
		assert.NoError(t, err)
		assert.Equal(t, movies, m)
	})

	t.Run("error while getting movies", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WillReturnError(errors.New(errTxt))

		_, err := processor.GetMovies(models.SortByRating)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting movies")
	})

	t.Run("error while getting actors", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(1, "name", "description", time.Time{}, 5))
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New(errTxt))

		_, err := processor.GetMovies(models.SortByRating)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting movie's actors")
		assert.Contains(t, err.Error(), "error while getting movies")
	})
}

func TestGetMoviesByActor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		actor := "name"
		movies := []models.MovieOut{
			{
				Id:          1,
				Name:        "name1",
				Description: "description1",
				ReleaseDate: time.Time{},
				Rating:      5,
				Actors:      []int{1},
			},
			{
				Id:          2,
				Name:        "name2",
				Description: "description2",
				ReleaseDate: time.Time{},
				Rating:      4,
				Actors:      []int{2},
			},
		}

		mock.ExpectQuery("SELECT").WithArgs(actor).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(movies[0].Id, movies[0].Name, movies[0].Description, movies[0].ReleaseDate, movies[0].Rating).AddRow(movies[1].Id, movies[1].Name, movies[1].Description, movies[1].ReleaseDate, movies[1].Rating))
		mock.ExpectQuery("SELECT").WithArgs(movies[0].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT").WithArgs(movies[1].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		m, err := processor.GetMoviesByActor(actor)
		assert.NoError(t, err)
		assert.Equal(t, movies, m)
	})

	t.Run("error while getting movies by actor", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		actor := "name"
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WithArgs(actor).WillReturnError(errors.New(errTxt))

		_, err := processor.GetMoviesByActor(actor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting movies by actor")
	})

	t.Run("error while getting actors", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		actor := "name"
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WithArgs(actor).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(1, "name", "description", time.Time{}, 5))
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New(errTxt))

		_, err := processor.GetMoviesByActor(actor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting movie's actors")
		assert.Contains(t, err.Error(), "error while getting movies by actor")
	})
}

func TestGetMoviesByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		name := "name"
		movies := []models.MovieOut{
			{
				Id:          1,
				Name:        "name1",
				Description: "description1",
				ReleaseDate: time.Time{},
				Rating:      5,
				Actors:      []int{1},
			},
			{
				Id:          2,
				Name:        "name2",
				Description: "description2",
				ReleaseDate: time.Time{},
				Rating:      4,
				Actors:      []int{2},
			},
		}

		mock.ExpectQuery("SELECT").WithArgs(name).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(movies[0].Id, movies[0].Name, movies[0].Description, movies[0].ReleaseDate, movies[0].Rating).AddRow(movies[1].Id, movies[1].Name, movies[1].Description, movies[1].ReleaseDate, movies[1].Rating))
		mock.ExpectQuery("SELECT").WithArgs(movies[0].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery("SELECT").WithArgs(movies[1].Id).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		m, err := processor.GetMoviesByName(name)
		assert.NoError(t, err)
		assert.Equal(t, movies, m)
	})

	t.Run("error while getting movies by name", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		name := "name"
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WithArgs(name).WillReturnError(errors.New(errTxt))

		_, err := processor.GetMoviesByName(name)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting movies by name")
	})

	t.Run("error while getting actors", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		name := "name"
		errTxt := "select error"
		mock.ExpectQuery("SELECT").WithArgs(name).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "release_date", "rating"}).AddRow(1, "name", "description", time.Time{}, 5))
		mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New(errTxt))

		_, err := processor.GetMoviesByName(name)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while getting movie's actors")
		assert.Contains(t, err.Error(), "error while getting movies by name")
	})
}

func TestUpdateActor(t *testing.T) {
	t.Run("success full", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		actor := models.ActorIn{
			Name:        "name",
			DateOfBirth: time.Now(),
			Gender:      "other",
		}
		id := 1

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Gender).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.DateOfBirth).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := processor.UpdateActor(id, actor)
		assert.NoError(t, err)
	})

	t.Run("success part 1", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		actor := models.ActorIn{
			Name:   "name",
			Gender: "other",
		}
		id := 1

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Gender).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := processor.UpdateActor(id, actor)
		assert.NoError(t, err)
	})

	t.Run("success part 2", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		actor := models.ActorIn{
			Name:        "name",
			DateOfBirth: time.Now(),
		}
		id := 1

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.DateOfBirth).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := processor.UpdateActor(id, actor)
		assert.NoError(t, err)
	})

	t.Run("error while starting transaction", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		actor := models.ActorIn{
			Name:        "name",
			Gender:      "other",
			DateOfBirth: time.Now(),
		}
		id := 1

		errTxt := "begin error"
		mock.ExpectBegin().WillReturnError(errors.New(errTxt))

		err := processor.UpdateActor(id, actor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errBeginTx.Error())
		assert.Contains(t, err.Error(), errTxt)
	})

	t.Run("error while updating actor 1", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		actor := models.ActorIn{

			Name:        "name",
			Gender:      "other",
			DateOfBirth: time.Now(),
		}
		id := 1

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Name).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateActor(id, actor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating actor")
	})

	t.Run("error while updating actor 2", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		actor := models.ActorIn{

			Name:        "name",
			Gender:      "other",
			DateOfBirth: time.Now(),
		}
		id := 1

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Gender).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateActor(id, actor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating actor")
	})

	t.Run("error while updating actor 3", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		actor := models.ActorIn{

			Name:        "name",
			Gender:      "other",
			DateOfBirth: time.Now(),
		}
		id := 1

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Name).WillReturnError(errors.New(errTxt))
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Gender).WillReturnError(errors.New(errTxt))
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.DateOfBirth).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateActor(id, actor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating actor")
	})

	t.Run("error while committing transaction", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		actor := models.ActorIn{
			Name: "name",
		}
		id := 1

		errTxt := "commit error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE actors").WithArgs(id, actor.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit().WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateActor(id, actor)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errCommitTx.Error())
		assert.Contains(t, err.Error(), errTxt)
	})
}

func TestUpdateMovie(t *testing.T) {
	t.Run("success full", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Description).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.ReleaseDate).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, *movie.Rating).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("DELETE FROM movie_actors").WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
		for _, a := range movie.Actors {
			mock.ExpectExec("INSERT INTO movie_actors").WithArgs(id, a).WillReturnResult(sqlmock.NewResult(0, 0))
		}
		mock.ExpectCommit()

		err := processor.UpdateMovie(id, movie)
		assert.NoError(t, err)
	})

	t.Run("success part 1", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Description).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.ReleaseDate).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, *movie.Rating).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("DELETE FROM movie_actors").WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
		for _, a := range movie.Actors {
			mock.ExpectExec("INSERT INTO movie_actors").WithArgs(id, a).WillReturnResult(sqlmock.NewResult(0, 0))
		}
		mock.ExpectCommit()

		err := processor.UpdateMovie(id, movie)
		assert.NoError(t, err)
	})

	t.Run("success part 2", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			ReleaseDate: time.Now(),
			Rating:      new(int),
		}
		id := 1
		*movie.Rating = 5

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.ReleaseDate).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, *movie.Rating).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := processor.UpdateMovie(id, movie)
		assert.NoError(t, err)
	})

	t.Run("error while starting transaction", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		errTxt := "begin error"
		mock.ExpectBegin().WillReturnError(errors.New(errTxt))

		err := processor.UpdateMovie(id, movie)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errBeginTx.Error())
		assert.Contains(t, err.Error(), errTxt)
	})

	t.Run("error while updating movie 1", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Name).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateMovie(id, movie)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating movie")
	})

	t.Run("error while updating movie 2", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Description).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateMovie(id, movie)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating movie")
	})

	t.Run("error while updating movie 3", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Name).WillReturnError(errors.New(errTxt))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Description).WillReturnError(errors.New(errTxt))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.ReleaseDate).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateMovie(id, movie)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating movie")
	})

	t.Run("error while updating movie 4", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Description).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.ReleaseDate).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, *movie.Rating).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateMovie(id, movie)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating movie")
	})

	t.Run("error while updating movie 5", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Description).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.ReleaseDate).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, *movie.Rating).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("DELETE FROM movie_actors").WithArgs(id).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateMovie(id, movie)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating movie")
	})

	t.Run("error while updating movie 6", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name:        "name",
			Description: "description",
			ReleaseDate: time.Now(),
			Rating:      new(int),
			Actors:      []int{1, 2, 3},
		}
		id := 1
		*movie.Rating = 5

		errTxt := "update error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Description).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.ReleaseDate).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("UPDATE movies").WithArgs(id, *movie.Rating).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("DELETE FROM movie_actors").WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("INSERT INTO movie_actors").WithArgs(id, movie.Actors[0]).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("INSERT INTO movie_actors").WithArgs(id, movie.Actors[1]).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateMovie(id, movie)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), "error while updating movie")
	})

	t.Run("error while committing transaction", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}

		movie := models.MovieIn{
			Name: "name",
		}
		id := 1

		errTxt := "commit error"
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE movies").WithArgs(id, movie.Name).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit().WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.UpdateMovie(id, movie)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errCommitTx.Error())
		assert.Contains(t, err.Error(), errTxt)
	})
}

func TestDeleteSmth(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		smth := true
		q := "smth"
		wrap := "error"

		mock.ExpectBegin()
		mock.ExpectExec(q).WithArgs(smth).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := processor.deleteSmth(q, wrap, smth)
		assert.NoError(t, err)
	})

	t.Run("error while starting transaction", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		smth := true
		q := "smth"
		wrap := "error"

		errTxt := "begin error"
		mock.ExpectBegin().WillReturnError(errors.New(errTxt))

		err := processor.deleteSmth(q, wrap, smth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errBeginTx.Error())
		assert.Contains(t, err.Error(), errTxt)
	})

	t.Run("error while deleting smth", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		smth := true
		q := "smth"
		wrap := "error"

		errTxt := "delete error"
		mock.ExpectBegin()
		mock.ExpectExec(q).WithArgs(smth).WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.deleteSmth(q, wrap, smth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errTxt)
		assert.Contains(t, err.Error(), wrap)
	})

	t.Run("error while committing transaction", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		smth := true
		q := "smth"
		wrap := "error"

		errTxt := "commit error"
		mock.ExpectBegin()
		mock.ExpectExec(q).WithArgs(smth).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit().WillReturnError(errors.New(errTxt))
		mock.ExpectRollback()

		err := processor.deleteSmth(q, wrap, smth)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errCommitTx.Error())
		assert.Contains(t, err.Error(), errTxt)
	})
}

func TestDeleteActor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		id := 1

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM actors").WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := processor.DeleteActor(id)
		assert.NoError(t, err)
	})
}

func TestDeleteMovie(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.Newx()
		defer db.Close()
		processor := dbProcessor{db: db}
		id := 1

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM movies").WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := processor.DeleteMovie(id)
		assert.NoError(t, err)
	})
}

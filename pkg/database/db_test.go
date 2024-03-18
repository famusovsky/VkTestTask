package database

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOpenViaDsn(t *testing.T) {
	t.Run("successful opening", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		var (
			dsn    = "test_dsn"
			driver = "postgres"
		)
		openMock := func(driverName, dataSourceName string) (*sql.DB, error) {
			assert.Equal(t, dsn, dataSourceName)
			assert.Equal(t, driver, driverName)
			return mockDB, nil
		}

		mock.ExpectPing()
		db, err := OpenViaDsn(dsn, driver, openMock)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty dsn", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		setEnv("admin", "qwerty", "127.0.0.1", "8080", "postgres")
		defer unsetEnv()

		openMock := func(driverName, dataSourceName string) (*sql.DB, error) {
			assert.Equal(t, "postgres://admin:qwerty@127.0.0.1:8080/postgres?sslmode=disable", dataSourceName)
			return mockDB, nil
		}

		mock.ExpectPing()
		db, err := OpenViaDsn("", "postgres", openMock)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error opening", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		testErr := errors.New("test")
		openMock := func(driverName, dataSourceName string) (*sql.DB, error) {
			return nil, testErr
		}

		db, err := OpenViaDsn("tmp", "tmp", openMock)
		assert.ErrorIs(t, err, testErr)
		assert.Nil(t, db)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOpenViaEnvVars(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	setEnv("admin", "qwerty", "127.0.0.1", "8080", "postgres")
	defer unsetEnv()

	openMock := func(driverName, dataSourceName string) (*sql.DB, error) {
		assert.Equal(t, "postgres://admin:qwerty@127.0.0.1:8080/postgres?sslmode=disable", dataSourceName)
		assert.Equal(t, "postgres", driverName)
		return mockDB, nil
	}

	mock.ExpectPing()
	db, err := OpenViaDsn("", "postgres", openMock)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDsnFromEnv(t *testing.T) {
	var (
		firstEnv  = []string{"user", "password", "localhost", "5432", "testdb"}
		secondEnv = []string{"admin", "qwerty", "127.0.0.1", "8080", "postgres"}

		firstExpected  = "postgres://user:password@localhost:5432/testdb?sslmode=disable"
		secondExpected = "postgres://admin:qwerty@127.0.0.1:8080/postgres?sslmode=disable"
	)
	t.Run("first case", func(t *testing.T) {
		setEnv(firstEnv[0], firstEnv[1], firstEnv[2], firstEnv[3], firstEnv[4])
		defer unsetEnv()
		assert.Equal(t, firstExpected, getDsnFromEnv())
	})

	t.Run("second case", func(t *testing.T) {
		setEnv(secondEnv[0], secondEnv[1], secondEnv[2], secondEnv[3], secondEnv[4])
		defer unsetEnv()
		assert.Equal(t, secondExpected, getDsnFromEnv())
	})
}

func unsetEnv() {
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
}

func setEnv(user, pswd, host, port, name string) {
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PASSWORD", pswd)
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port)
	os.Setenv("DB_NAME", name)
}

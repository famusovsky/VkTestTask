// Пакет для работы с БД
package database

import (
	"database/sql"
	"fmt"
	"os"
)

// OpenViaEnvVars - открытие БД через переменные окружения.
// Возвращает БД и ошибку.
func OpenViaEnvVars(driver string, openFunc func(string, string) (*sql.DB, error)) (*sql.DB, error) {
	return OpenViaDsn(getDsnFromEnv(), driver, openFunc)
}

// OpenViaDsn - открытие БД через строку DSN.
// Принимает строку DSN.
// Возвращает БД и ошибку.
func OpenViaDsn(dsn, driver string, openFunc func(string, string) (*sql.DB, error)) (*sql.DB, error) {
	if dsn == "" {
		dsn = getDsnFromEnv()
	}

	db, err := openFunc(driver, dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// getDsnFromEnv - получение строки DSN из переменных окружения.
func getDsnFromEnv() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	return dsn
}

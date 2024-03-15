package postgres

import "database/sql"

// TODO
type DbHandler interface {
}

func GetHandler(db *sql.DB, createTables bool) (DbHandler, error) {
	return dbProcessor{}, nil
}

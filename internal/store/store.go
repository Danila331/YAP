package store

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func ConnectToDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "../internal/store/data.db")
	if err != nil {
		return &sql.DB{}, err
	}

	return db, nil
}

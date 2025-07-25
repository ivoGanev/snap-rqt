package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
)


type SQLiteDb struct {
	db *sql.DB
}

func NewDb(path string) (*SQLiteDb, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	database := &SQLiteDb{db: db}


	return database, nil
}


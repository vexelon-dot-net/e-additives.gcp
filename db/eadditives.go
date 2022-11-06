package db

import (
	"database/sql"
	"fmt"
)

const dateTimeLayout = "2006-01-02 15:04:05"

type Row interface {
	Scan(...interface{}) error
}

type RowScanner interface {
	ScanFrom(Row) error
}

var db *sql.DB

func Open(path string) (err error) {
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("Error opening database at %s : %w", path, err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("Failed database ping : %w", err)
	}

	return nil
}

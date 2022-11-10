package db

import (
	"database/sql"
	"fmt"
)

const dateTimeLayout = "2006-01-02 15:04:05"

type DBChannel struct {
	db *sql.DB
}

type DBProvider struct {
	Locales   localesChannel
	Additives additivesChannel
}

type Row interface {
	Scan(...interface{}) error
}

type RowScanner interface {
	ScanFrom(Row) error
}

func NewProvider(path string) (provider *DBProvider, err error) {
	sqlDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Error opening database at %s : %w", path, err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("Failed database ping : %w", err)
	}

	provider = &DBProvider{
		Locales:   localesChannel{DBChannel{sqlDB}},
		Additives: additivesChannel{DBChannel{sqlDB}},
	}

	return provider, nil
}

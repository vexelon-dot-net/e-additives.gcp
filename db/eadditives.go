package db

import (
	"database/sql"
	"fmt"
)

type Locale struct {
	Id      int    `json:"id"`
	Code    string `json:"code"`
	Enabled bool   `json:"enabled"`
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

func FetchAllLocales() ([]*Locale, error) {
	rows, err := db.Query(`SELECT * FROM ead_Locale`)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all locales: %w", err)
	}
	defer rows.Close()

	locales := make([]*Locale, 0)
	for rows.Next() {
		device := new(Locale)
		err := rows.Scan(&device.Id, &device.Code, &device.Enabled)
		if err != nil {
			return nil, fmt.Errorf("Error scanning locale row: %w", err)
		}
		locales = append(locales, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return locales, nil
}

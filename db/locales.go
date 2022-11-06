package db

import "fmt"

type Locale struct {
	Id      int    `json:"id"`
	Code    string `json:"code"`
	Enabled bool   `json:"enabled"`
}

func FetchAllLocales() ([]*Locale, error) {
	rows, err := db.Query(`SELECT * FROM ead_Locale`)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all locales: %w", err)
	}
	defer rows.Close()

	locales := make([]*Locale, 0)
	for rows.Next() {
		locale := new(Locale)
		err := rows.Scan(&locale.Id, &locale.Code, &locale.Enabled)
		if err != nil {
			return nil, fmt.Errorf("Error scanning locale row: %w", err)
		}
		locales = append(locales, locale)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return locales, nil
}

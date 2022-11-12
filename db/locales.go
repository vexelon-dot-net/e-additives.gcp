package db

import (
	"fmt"
)

type localesChannel struct {
	DBChannel
}

type Locale struct {
	Id      int    `json:"id"`
	Code    string `json:"code"`
	Enabled bool   `json:"enabled"`
	Url     string `json:"url,omitempty"`
}

func (loc *Locale) ScanFrom(r Row) (err error) {
	if err = r.Scan(&loc.Id, &loc.Code, &loc.Enabled); err != nil {
		return fmt.Errorf("Error scanning locale row: %w", err)
	}
	return err
}

func (chn *localesChannel) FetchAll() ([]*Locale, error) {
	rows, err := chn.db.Query(`SELECT * FROM ead_Locale`)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all locales: %w", err)
	}
	defer rows.Close()

	locales := make([]*Locale, 0)
	for rows.Next() {
		loc := new(Locale)
		if err := loc.ScanFrom(rows); err != nil {
			return nil, fmt.Errorf("Error scanning locale row: %w", err)
		}
		locales = append(locales, loc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return locales, nil
}

func (chn *localesChannel) FetchOne(locId int) (*Locale, error) {
	loc := new(Locale)

	err := loc.ScanFrom(chn.db.QueryRow(`SELECT * FROM ead_Locale WHERE id=$1`, locId))
	if err != nil {
		return nil, fmt.Errorf("Error fetching single locale '%d': %w", locId, err)
	}

	return loc, nil
}

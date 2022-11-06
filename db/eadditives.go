package db

import (
	"database/sql"
	"fmt"
	"time"
)

const DATE_TIME_LAYOUT = "2006-01-02 15:04:05"

type Locale struct {
	Id      int    `json:"id"`
	Code    string `json:"code"`
	Enabled bool   `json:"enabled"`
}

type AdditiveCategory struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	LastUpdate       string    `json:"last_update"`
	LastUpdateParsed time.Time `json:"-"`
	Additives        int       `json:"additives"`
}

// type AdditiveCategory struct {
// 	Id               int       `json:"id"`
// 	Category         int       `json:"category"`
// 	LastUpdate       string    `json:"last_update"`
// 	LastUpdateParsed time.Time `json:"-"`
// }

// type AdditiveCategoryProps struct {
// 	Id               int       `json:"id"`
// 	CategoryId       int       `json:"-"`
// 	LocaleId         int       `json:"-"`
// 	Name             string    `json:"name"`
// 	Description      string    `json:"description"`
// 	LastUpdate       string    `json:"last_update"`
// 	LastUpdateParsed time.Time `json:"-"`
// }

type Row interface {
	Scan(...interface{}) error
}

type RowScanner interface {
	ScanRow(Row) error
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

func (ac *AdditiveCategory) ScanRow(r Row) (err error) {
	if err = r.Scan(&ac.Id, &ac.Name, &ac.Description, &ac.LastUpdate, &ac.Additives); err != nil {
		return fmt.Errorf("Error scanning AdditiveCategory row: %w", err)
	}
	ac.LastUpdateParsed, err = time.Parse(DATE_TIME_LAYOUT, ac.LastUpdate)
	return err
}

func FetchAllCategories(loc Locale) ([]*AdditiveCategory, error) {
	rows, err := db.Query(`
		SELECT c.id, p.name, p.description, p.last_update,
		(SELECT COUNT(id) FROM ead_Additive as a WHERE a.category_id=c.id) as additives
		FROM ead_AdditiveCategory as c
		LEFT JOIN ead_AdditiveCategoryProps as p ON p.category_id = c.id
		WHERE p.locale_id = $1
	`, loc.Id)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all categories: %w", err)
	}
	defer rows.Close()

	cats := make([]*AdditiveCategory, 0)
	for rows.Next() {
		cat := new(AdditiveCategory)
		if err = cat.ScanRow(rows); err != nil {
			return nil, fmt.Errorf("Error scanning locale row: %w", err)
		}
		cats = append(cats, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cats, nil
}

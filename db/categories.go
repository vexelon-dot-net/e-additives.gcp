package db

import (
	"fmt"
	"time"
)

type categoriesChannel struct {
	DBChannel
}

type AdditiveCategory struct {
	Id               int       `json:"-"`
	Category         int       `json:"category"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	LastUpdate       string    `json:"last_update"`
	LastUpdateParsed time.Time `json:"-"`
	Additives        int       `json:"additives"`
	Url              string    `json:"url,omitempty"`
}

func (ac *AdditiveCategory) ScanFrom(r Row) (err error) {
	if err = r.Scan(&ac.Id, &ac.Category, &ac.Name, &ac.Description,
		&ac.LastUpdate, &ac.Additives); err != nil {
		return fmt.Errorf("Error scanning additive category row: %w", err)
	}
	ac.LastUpdateParsed, err = time.Parse(dateTimeLayout, ac.LastUpdate)
	return err
}

func (chn *categoriesChannel) FetchAll(loc Locale) ([]*AdditiveCategory, error) {
	rows, err := chn.db.Query(`
		SELECT c.id, c.category, p.name, p.description, p.last_update,
		(SELECT COUNT(id) FROM ead_Additive AS a WHERE a.category_id=c.id) AS additives
		FROM ead_AdditiveCategory AS c
		LEFT JOIN ead_AdditiveCategoryProps AS p ON p.category_id = c.id
		WHERE p.locale_id = $1
	`, loc.Id)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all categories: %w", err)
	}
	defer rows.Close()

	cats := make([]*AdditiveCategory, 0)
	for rows.Next() {
		cat := new(AdditiveCategory)
		if err = cat.ScanFrom(rows); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cats, nil
}

func (chn *categoriesChannel) FetchOne(category int, loc Locale) (*AdditiveCategory, error) {
	cat := new(AdditiveCategory)

	err := cat.ScanFrom(chn.db.QueryRow(`
		SELECT c.id, c.category, p.name, p.description, p.last_update,
		(SELECT COUNT(id) FROM ead_Additive AS a WHERE a.category_id=c.id) AS additives
		FROM ead_AdditiveCategory AS c
		LEFT JOIN ead_AdditiveCategoryProps AS p ON p.category_id = c.id
		WHERE c.category = $1 AND p.locale_id = $2
	`, category, loc.Id))
	if err != nil {
		return nil, fmt.Errorf("Error fetching single category '%d': %w", category, err)
	}

	return cat, nil
}

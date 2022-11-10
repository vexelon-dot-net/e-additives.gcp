package db

import (
	"fmt"
	"time"
)

type categoriesChannel struct {
	DBChannel
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

func (ac *AdditiveCategory) ScanFrom(r Row) (err error) {
	if err = r.Scan(&ac.Id, &ac.Name, &ac.Description, &ac.LastUpdate, &ac.Additives); err != nil {
		return fmt.Errorf("Error scanning AdditiveCategory row: %w", err)
	}
	ac.LastUpdateParsed, err = time.Parse(dateTimeLayout, ac.LastUpdate)
	return err
}

func (chn *categoriesChannel) FetchAll(loc Locale) ([]*AdditiveCategory, error) {
	rows, err := chn.db.Query(`
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

func (chn *categoriesChannel) FetchOne(catId int, loc Locale) (*AdditiveCategory, error) {
	cat := new(AdditiveCategory)

	err := cat.ScanFrom(chn.db.QueryRow(`
		SELECT c.id, p.name, p.description, p.last_update,
		(SELECT COUNT(id) FROM ead_Additive as a WHERE a.category_id=c.id) as additives
		FROM ead_AdditiveCategory as c
		LEFT JOIN ead_AdditiveCategoryProps as p ON p.category_id = c.id
		WHERE c.id = $1 AND p.locale_id = $2
	`, catId, loc.Id))
	if err != nil {
		return nil, fmt.Errorf("Error fetching single category '%d': %w", catId, err)
	}

	return cat, nil
}

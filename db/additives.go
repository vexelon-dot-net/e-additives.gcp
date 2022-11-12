package db

import (
	"database/sql"
	"fmt"
	"time"
)

type additivesChannel struct {
	DBChannel
	Categories categoriesChannel
}

type AdditiveMeta struct {
	Id               int       `json:"-"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	LastUpdate       string    `json:"last_update"`
	LastUpdateParsed time.Time `json:"-"`
}

type Additive struct {
	Id               int       `json:"-"`
	CatId            int       `json:"category_id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	Veg              string    `json:"veg"`
	Function         string    `json:"function"`
	Foods            string    `json:"foods"`
	Notice           string    `json:"notice"`
	Info             string    `json:"info"`
	LastUpdate       string    `json:"last_update"`
	LastUpdateParsed time.Time `json:"-"`
}

func (am *AdditiveMeta) ScanFrom(r Row) (err error) {
	if err = r.Scan(&am.Id, &am.Code, &am.LastUpdate, &am.Name); err != nil {
		return fmt.Errorf("Error scanning additive meta row: %w", err)
	}
	am.LastUpdateParsed, err = time.Parse(dateTimeLayout, am.LastUpdate)
	return err
}

func (a *Additive) ScanFrom(r Row) (err error) {
	var (
		veg    sql.NullBool
		status sql.NullString
		foods  sql.NullString
		notice sql.NullString
	)
	if err = r.Scan(&a.Id, &a.Code, &a.LastUpdate, &a.CatId, &a.Name, &status,
		&veg, &a.Function, &foods, &notice, &a.Info); err != nil {
		return fmt.Errorf("Error scanning additive row: %w", err)
	}
	a.Veg = yesNoEmptyIfNull(veg)
	a.Status = emptyIfNull(status)
	a.Foods = emptyIfNull(foods)
	a.Notice = emptyIfNull(notice)
	a.LastUpdateParsed, err = time.Parse(dateTimeLayout, a.LastUpdate)
	return err
}

func (chn *additivesChannel) FetchAll(loc Locale) ([]*AdditiveMeta, error) {
	rows, err := chn.db.Query(`
		SELECT a.id, a.code, a.last_update,
		(SELECT value_str FROM ead_AdditiveProps WHERE additive_id = a.id 
			AND key_name = 'name' AND locale_id = $1) AS name
		FROM ead_Additive AS a
	`, loc.Id)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all additives: %w", err)
	}
	defer rows.Close()

	additives := make([]*AdditiveMeta, 0)
	for rows.Next() {
		am := new(AdditiveMeta)
		if err = am.ScanFrom(rows); err != nil {
			return nil, err
		}
		additives = append(additives, am)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return additives, nil
}

func (chn *additivesChannel) FetchOne(code string, loc Locale) (*Additive, error) {
	a := new(Additive)

	err := a.ScanFrom(chn.db.QueryRow(`
		SELECT a.id, a.code, a.last_update, a.category_id,
		(SELECT value_str FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'name' AND locale_id = $1) AS name,
		(SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'status' AND locale_id = $1) AS status,
		(SELECT value_str FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'veg' AND locale_id = $1) AS veg,
		(SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'function' AND locale_id = $1) AS function,
		(SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'foods' AND locale_id = $1) AS foods,
		(SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'notice' AND locale_id = $1) AS notice,
		(SELECT value_big_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'info' AND locale_id = $1) AS info
		FROM ead_Additive AS a 
		WHERE a.code = $2
	`, loc.Id, code))
	if err != nil {
		return nil, fmt.Errorf("Error fetching single additive '%s': %w", code, err)
	}

	return a, nil
}

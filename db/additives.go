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
	Name             string    `json:"name,omitempty"`
	LastUpdate       string    `json:"last_update,omitempty"`
	LastUpdateParsed time.Time `json:"-"`
	Url              string    `json:"url,omitempty"`
}

type Additive struct {
	Id               int       `json:"-"`
	Category         int       `json:"category"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	Status           string    `json:"status,omitempty"`
	Veg              string    `json:"veg,omitempty"`
	Function         string    `json:"function,omitempty"`
	Foods            string    `json:"foods,omitempty"`
	Notice           string    `json:"notice,omitempty"`
	Info             string    `json:"info,omitempty"`
	LastUpdate       string    `json:"last_update,omitempty"`
	LastUpdateParsed time.Time `json:"-"`
	Url              string    `json:"url,omitempty"`
}

func (am *AdditiveMeta) ScanFrom(r Row) (err error) {
	var name sql.NullString
	var lastUpdate sql.NullString
	if err = r.Scan(&am.Id, &am.Code, &lastUpdate, &name); err != nil {
		return fmt.Errorf("Error scanning additive meta row: %w", err)
	}
	am.Name = emptyIfNull(name)
	if lastUpdate.Valid {
		am.LastUpdate = lastUpdate.String
		am.LastUpdateParsed, err = time.Parse(dateTimeLayout, am.LastUpdate)
	}
	return err
}

func (a *Additive) ScanFrom(r Row) (err error) {
	var (
		name     sql.NullString
		status   sql.NullString
		veg      sql.NullBool
		function sql.NullString
		foods    sql.NullString
		notice   sql.NullString
		info     sql.NullString
	)
	if err = r.Scan(&a.Id, &a.Code, &a.LastUpdate, &a.Category, &name, &status,
		&veg, &function, &foods, &notice, &info); err != nil {
		return fmt.Errorf("Error scanning additive row: %w", err)
	}
	a.Name = emptyIfNull(name)
	a.Status = emptyIfNull(status)
	a.Veg = yesNoEmptyIfNull(veg)
	a.Function = emptyIfNull(function)
	a.Foods = emptyIfNull(foods)
	a.Notice = emptyIfNull(notice)
	a.Info = emptyIfNull(info)
	a.LastUpdateParsed, err = time.Parse(dateTimeLayout, a.LastUpdate)
	return err
}

func additiveRowsToArray(rows *sql.Rows) ([]*AdditiveMeta, error) {
	additives := make([]*AdditiveMeta, 0)
	for rows.Next() {
		am := new(AdditiveMeta)
		if err := am.ScanFrom(rows); err != nil {
			return nil, err
		}
		additives = append(additives, am)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return additives, nil
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

	return additiveRowsToArray(rows)
}

func (chn *additivesChannel) FetchAllByCategory(category int, loc Locale) ([]*AdditiveMeta, error) {
	rows, err := chn.db.Query(`
		SELECT a.id, a.code, a.last_update,
		(SELECT value_str FROM ead_AdditiveProps WHERE additive_id = a.id 
			AND key_name = 'name' AND locale_id = $1) AS name
		FROM ead_Additive AS a
		JOIN ead_AdditiveCategory AS c ON c.id=a.category_id
		WHERE c.category = $2
	`, loc.Id, category)
	if err != nil {
		return nil, fmt.Errorf("Error fetching category (%d) additives: %w", category, err)
	}
	defer rows.Close()

	return additiveRowsToArray(rows)
}

func (chn *additivesChannel) Search(keyword string) ([]*AdditiveMeta, error) {
	rows, err := chn.db.Query(`
		SELECT id, code, NULL, name FROM ead_AdditiveFTSI 
		WHERE ead_AdditiveFTSI MATCH $1 ORDER BY rank`,
		keyword)
	if err != nil {
		return nil, fmt.Errorf("Error fetching additives using keyword (%s) additives: %w", keyword, err)
	}
	defer rows.Close()

	return additiveRowsToArray(rows)
}

func (chn *additivesChannel) FetchOne(code string, loc Locale) (*Additive, error) {
	a := new(Additive)

	err := a.ScanFrom(chn.db.QueryRow(`
		SELECT a.id, a.code, a.last_update, c.category,
		(SELECT value_str FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'name' AND locale_id = $1) AS name,
		(SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'status' AND locale_id = $1) AS status,
		(SELECT value_str FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'veg' AND locale_id = $1) AS veg,
		(SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'function' AND locale_id = $1) AS function,
		(SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'foods' AND locale_id = $1) AS foods,
		(SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'notice' AND locale_id = $1) AS notice,
		(SELECT value_big_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'info' AND locale_id = $1) AS info
		FROM ead_Additive AS a 
		JOIN ead_AdditiveCategory AS c ON c.id=a.category_id
		WHERE a.code = $2
	`, loc.Id, code))
	if err != nil {
		return nil, fmt.Errorf("Error fetching single additive '%s': %w", code, err)
	}

	return a, nil
}

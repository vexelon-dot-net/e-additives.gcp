//go:build eadfts
// +build eadfts

package db

import "fmt"

func (chn *additivesChannel) Search(keyword string, loc Locale) ([]*AdditiveMeta, error) {
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

//go:build !eadfts
// +build !eadfts

package db

import "fmt"

func (chn *additivesChannel) Search(keyword string, loc Locale) ([]*AdditiveMeta, error) {
	rows, err := chn.db.Query(`
		SELECT p.additive_id AS id, a.code, NULL,
			(SELECT value_str FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'name' AND locale_id = $1) AS name
		FROM ead_Additive AS a
		LEFT JOIN ead_AdditiveProps AS p ON p.additive_id = a.id
		WHERE p.locale_id = $1 AND ((a.code LIKE $2 || '%') OR 
			(p.key_name = 'name' AND p.value_str LIKE '%' || $2 || '%') OR 
			(p.key_name = 'status' AND p.value_text LIKE '%' || $2 || '%') OR 
			(p.key_name = 'function' AND p.value_text LIKE '%' || $2 || '%') OR 
			(p.key_name = 'foods' AND p.value_text LIKE '%' || $2 || '%') OR 
			(p.key_name = 'notice' AND p.value_text LIKE '%' || $2 || '%') OR 
			(p.key_name = 'info' AND p.value_text LIKE '%' || $2 || '%')
		)
		GROUP BY p.additive_id, a.code, NULL, name
		`, loc.Id, keyword)
	if err != nil {
		return nil, fmt.Errorf("Error fetching additives using keyword (%s) additives: %w", keyword, err)
	}
	defer rows.Close()

	return additiveRowsToArray(rows)
}

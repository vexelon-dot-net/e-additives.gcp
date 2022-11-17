package db

import (
	"fmt"
)

type apiKeysChannel struct {
	DBChannel
}

type ApiKey struct {
	Id      int
	Key     string
	Host    string
	Enabled bool
}

func (apiKey *ApiKey) ScanFrom(r Row) (err error) {
	if err = r.Scan(&apiKey.Id, &apiKey.Key, &apiKey.Host, &apiKey.Enabled); err != nil {
		return fmt.Errorf("Error scanning api key row: %w", err)
	}
	return err
}

func (chn *apiKeysChannel) FetchAll() ([]*ApiKey, error) {
	rows, err := chn.db.Query(`SELECT * FROM ead_ApiKey`)
	if err != nil {
		return nil, fmt.Errorf("Error fetching all api keys: %w", err)
	}
	defer rows.Close()

	apiKeys := make([]*ApiKey, 0)
	for rows.Next() {
		apiKey := new(ApiKey)
		if err := apiKey.ScanFrom(rows); err != nil {
			return nil, fmt.Errorf("Error scanning api key row: %w", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (chn *apiKeysChannel) FetchOne(key string) (*ApiKey, error) {
	apiKey := new(ApiKey)

	err := apiKey.ScanFrom(chn.db.QueryRow(`SELECT * FROM ead_ApiKey WHERE key='$1'`, key))
	if err != nil {
		return nil, fmt.Errorf("Error fetching single api key '%s': %w", key, err)
	}

	return apiKey, nil
}

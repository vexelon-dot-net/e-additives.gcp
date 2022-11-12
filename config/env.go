package config

import (
	"fmt"
	"os"
	"strconv"
)

func CreateFromEnv() (config *Config, err error) {
	config = new(Config)

	value, isPresent := os.LookupEnv("PORT")
	if !isPresent {
		return nil, fmt.Errorf("Cannot find env `PORT`!")
	}
	config.ListenPort, _ = strconv.Atoi(value)

	config.ListenAddress, isPresent = os.LookupEnv("HOST")
	if !isPresent {
		config.ListenAddress = "" // default IPv6
	}

	config.DatabasePath, isPresent = os.LookupEnv("DB_PATH")
	if !isPresent {
		return nil, fmt.Errorf("Cannot find env `DB_PATH`!")
	}

	if err = verifyPath(config.DatabasePath, "DatabasePath", true); err != nil {
		return nil, err
	}

	value, isPresent = os.LookupEnv("DEVMODE")
	if isPresent {
		if config.IsDevMode, err = strconv.ParseBool(value); err != nil {
			return nil, fmt.Errorf("Cannot parse env `DEVMODE`! %w", err)
		}
	} else {
		config.IsDevMode = false
	}

	return config, nil
}

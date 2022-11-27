package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// server listen address
	ListenAddress string
	// server listen port
	ListenPort int
	// path to the sqlite database
	DatabasePath string
	// enbale dev mode behavior
	IsDevMode bool
}

func verifyPath(path string, what string, mustExist bool) error {
	if len(path) < 1 {
		return fmt.Errorf("Error: %s path not specified!\n\n", what)
	}

	if mustExist {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("Error: %s path not found at '%s'!\n\n", what, path)
		}
	}

	return nil
}

func NewFromEnv() (config *Config, err error) {
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

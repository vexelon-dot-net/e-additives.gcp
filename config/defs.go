package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	// server listen port
	ListenPort int
	// path to the sqlite database
	DatabasePath string
)

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

func ParseConfig() error {
	value, isPresent := os.LookupEnv("PORT")
	if !isPresent {
		return fmt.Errorf("Cannot find env `PORT`!")
	}
	ListenPort, _ = strconv.Atoi(value)

	DatabasePath, isPresent = os.LookupEnv("DB_PATH")
	if !isPresent {
		return fmt.Errorf("Cannot find env `DB_PATH`!")
	}

	if err := verifyPath(DatabasePath, "DatabasePath", true); err != nil {
		return err
	}

	return nil
}

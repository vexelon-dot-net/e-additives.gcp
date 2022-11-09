package config

import (
	"fmt"
	"os"
)

type Config struct {
	// server listen address
	ListenAddress string
	// server listen port
	ListenPort int
	// path to the sqlite database
	DatabasePath string
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

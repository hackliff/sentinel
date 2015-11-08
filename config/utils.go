package config

import (
	"fmt"
	"os"
)

// NOTE maybe fmt.Sprintf("%s.yml", agent.nodename)
// sentinel will try to read this file if environment variable
// SENTINEL_CONF_PATH isn't set
const DEFAULT_CONF_PATH = "./sentinel.yml"

func missingFile(filePath string) bool {
	_, err := os.Stat(filePath)
	os.IsNotExist(err)
	return os.IsNotExist(err)
}

func Path() (string, error) {
	filePath := os.Getenv("SENTINEL_CONF_PATH")
	if filePath == "" {
		filePath = DEFAULT_CONF_PATH
	}

	if missingFile(filePath) {
		return "", fmt.Errorf("configuration file %s not found", filePath)
	}

	return filePath, nil
}

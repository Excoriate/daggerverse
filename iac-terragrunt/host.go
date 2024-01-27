package main

import (
	"errors"
	"os"
	"path/filepath"
)

func toAbsolutePathAndValidate(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Check if the path exists.
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return "", errors.New("path does not exist")
	}
	if err != nil {
		return "", err
	}

	// Check if the path is a directory.
	if !info.IsDir() {
		return "", errors.New("path is not a directory")
	}

	return absPath, nil
}

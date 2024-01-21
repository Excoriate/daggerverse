package utils

import (
	"os"
	"path/filepath"
)

func ConvertToAbs(src string) (string, error) {
	absPath, err := filepath.Abs(src)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func DirExist(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return err
	}

	return nil
}

func DirHasFilesWithExtension(path, extension string) error {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == extension {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return err
	}

	return nil
}

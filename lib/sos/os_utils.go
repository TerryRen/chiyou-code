package sos

import (
	"errors"
	"os"
	"path/filepath"
)

// Create file and create directory if not exists
func CreateFile(name string) (fp *os.File, err error) {
	folder := filepath.Dir(name)
	if err = CreateFolder(folder); err != nil {
		return nil, err
	}
	return os.Create(name)
}

// Create directory if not exists else do nothing
func CreateFolder(dir string) (err error) {
	// Check if the directory exists
	_, err = os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		// Create the directory if it does not exist
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return err
}

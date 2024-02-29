package sos

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
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

// Create directory if not exists else do nothing
func CreateFolderWithModel(dir string, perm fs.FileMode) (err error) {
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

func CopyFolder(src, dst string) (err error) {
	if strings.EqualFold(src, dst) {
		return nil
	}
	// src walk
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) (oerr error) {
		if err != nil {
			return err
		}
		logrus.Info(path)
		outpath := filepath.Join(dst, strings.TrimPrefix(path, src))
		if info.IsDir() {
			return CreateFolderWithModel(outpath, info.Mode())
		} else {
			// Check the file mode of the entry
			switch info.Mode() & os.ModeType {
			case os.ModeSymlink:
				if oerr = CopySymLink(path, outpath); oerr != nil {
					return oerr
				}
			default:
				if oerr = CopyFileWithModel(path, outpath, info.Mode()); oerr != nil {
					return oerr
				}
			}
		}
		return nil
	})
}

func CopyFileWithModel(src, dst string, perm fs.FileMode) (err error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, perm)
}

func CopySymLink(src, dst string) (err error) {
	// Get the link target of the source link
	linkTarget, err := os.Readlink(src)
	if err != nil {
		return fmt.Errorf("couldn't read source link: %s", err)
	}
	// Create the destination link with the same link target as the source link
	err = os.Symlink(linkTarget, dst)
	if err != nil {
		return fmt.Errorf("couldn't create dest link: %s", err)
	}
	return nil
}

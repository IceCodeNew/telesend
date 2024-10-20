package fsHelper

import (
	"fmt"
	"os"
	"path/filepath"
)

func CloseThenRemove(file *os.File) error {
	// it does not matter if the close function fails
	_ = file.Close()
	return os.Remove(file.Name())
}

func CreateDir(path string) error {
	path = filepath.Dir(path)

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(path, ".write-test-")
	if err != nil {
		return err
	}
	defer CloseThenRemove(tmpFile)

	return nil
}

func IsRegularFile(path string) (bool, error) {
	if _info, err := os.Stat(path); err != nil {
		return false, err
	} else if _info.Mode().IsRegular() {
		return true, nil
	}
	return false, fmt.Errorf("%s is not a regular file", path)
}

package utils

import (
	"os"
)

// CreateDirectory creates a single directory if it doesn't exist
func CreateDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

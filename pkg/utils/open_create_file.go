package utils

import (
	"os"
	"path"
)

func OpenOrCreateFile(filename string) (*os.File, error) {
	dirname := path.Dir(filename)
	if err := os.MkdirAll(dirname, 0666); err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

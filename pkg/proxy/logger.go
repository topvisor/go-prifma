package proxy

import (
	"log"
	"os"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

func (t *Logger) SetLogger(logger *log.Logger) error {
	if err := t.Close(); err != nil {
		return err
	}

	t.logger = logger

	return nil
}

func (t *Logger) SetFile(filename string) error {
	var err error

	if err = t.Close(); err != nil {
		return err
	}
	if t.file, err = os.Open(filename); err != nil {
		return err
	}

	t.logger = log.New(t.file, "", log.LstdFlags)

	return nil
}

func (t *Logger) Close() error {
	if t.file != nil {
		if err := t.file.Close(); err != nil {
			return err
		}
	}

	t.logger = nil
	t.file = nil

	return nil
}

func (t *Logger) Println(v ...interface{}) {
	if t.logger != nil {
		t.logger.Println(v...)
	}
}

func (t *Logger) Fatalln(v ...interface{}) {
	t.Println(v...)
	os.Exit(1)
}

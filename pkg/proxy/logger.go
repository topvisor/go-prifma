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

func (t *Logger) IsInited() bool {
	return t.logger != nil
}

func (t *Logger) Println(v ...interface{}) {
	if t.IsInited() {
		t.logger.Println(v...)
	}
}

func (t *Logger) Printf(format string, v ...interface{}) {
	if t.IsInited() {
		t.logger.Printf(format, v...)
	}
}

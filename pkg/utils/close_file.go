package utils

import (
	"io"
	"os"
	"time"
)

const CloseFileTimeout = time.Second * 2

func CloseFile(closer io.Closer) {
	if err := closer.Close(); err != nil && err != os.ErrClosed && err != os.ErrNotExist {
		time.Sleep(CloseFileTimeout)
		_ = closer.Close()
	}
}

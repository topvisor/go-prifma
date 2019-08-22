package proxy

import (
	"io"
	"os"
	"time"
)

const (
	closeTimeout = time.Second * 2
)

func transfer(src io.ReadCloser, dst io.WriteCloser) {
	_, _ = io.Copy(dst, src)

	closeFile(dst)
	closeFile(src)
}

func closeFile(closer io.Closer) {
	if err := closer.Close(); err != nil && err != os.ErrClosed && err != os.ErrNotExist {
		time.Sleep(closeTimeout)
		_ = closer.Close()
	}
}

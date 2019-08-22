package proxy

import (
	"context"
	"io"
	"os"
	"time"
)

const closeTimeout = time.Second * 2

func contextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	} else {
		return context.WithCancel(ctx)
	}
}

func closeFile(closer io.Closer) {
	if err := closer.Close(); err != nil && err != os.ErrClosed && err != os.ErrNotExist {
		time.Sleep(closeTimeout)
		_ = closer.Close()
	}
}

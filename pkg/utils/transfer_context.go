package utils

import (
	"context"
	"io"
	"time"
)

const BufferSize = 1024 * 32

func Transfer(readTimeout time.Duration, writeTimeout time.Duration, src io.ReadCloser, dst io.WriteCloser) {
	var nr, nw int
	var err error
	var ctx context.Context
	var cancel context.CancelFunc

	buf := make([]byte, BufferSize)

	for {
		if readTimeout == 0 {
			ctx, cancel = context.WithCancel(context.Background())
		} else {
			ctx, cancel = context.WithTimeout(context.Background(), readTimeout)
		}

		go func() {
			nr, err = src.Read(buf)
			cancel()
		}()

		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded || err != nil {
			break
		}

		if nr > 0 {
			if writeTimeout == 0 {
				ctx, cancel = context.WithCancel(context.Background())
			} else {
				ctx, cancel = context.WithTimeout(context.Background(), writeTimeout)
			}

			go func() {
				nw, err = dst.Write(buf[:nr])
				cancel()
			}()

			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded || err != nil || nr != nw {
				break
			}
		}
	}

	CloseFile(dst)
	CloseFile(src)
}

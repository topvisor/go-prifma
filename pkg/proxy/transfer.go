package proxy

import (
	"context"
	"io"
)

const (
	bufferSize      = 1024 * 16
	maxBuffersCount = 4
)

func transfer(src io.ReadCloser, dst io.WriteCloser) {
	defer src.Close()
	defer dst.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	buffersChan := make(chan []byte, maxBuffersCount)

	go func() {
		defer cancel()
		readToChannelContext(src, buffersChan, ctx)
	}()

	go func() {
		defer cancel()
		writeFromChannel(dst, buffersChan)
	}()

	<-ctx.Done()
}

func readToChannelContext(r io.Reader, ch chan []byte, ctx context.Context) {
	defer close(ch)

	buffer := make([]byte, bufferSize)

	for done := false; !done; {
		if nr, err := r.Read(buffer); err == nil {
			select {
			case <-ctx.Done():
				done = true
			case ch <- buffer[:nr]:
			}
		} else {
			done = true
		}
	}
}

func writeFromChannel(w io.Writer, ch chan []byte) {
	for buffer := range ch {
		if nw, err := w.Write(buffer); err != nil || nw != len(buffer) {
			break
		}
	}
}

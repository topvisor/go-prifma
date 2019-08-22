package proxy

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

const bufferSize = 1024 * 32

type responseWriterTunnel struct {
	DestConn     net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	clientConn net.Conn
	code       int
}

func (t *responseWriterTunnel) GetCode() int {
	return t.code
}

func (t *responseWriterTunnel) Write(rw http.ResponseWriter) error {
	rw.WriteHeader(http.StatusOK)
	t.code = http.StatusOK

	clientConn, _, hijackError := rw.(http.Hijacker).Hijack()
	if hijackError != nil {
		if err := t.DestConn.Close(); err != nil {
			_ = t.DestConn.Close()
		}

		http.Error(rw, hijackError.Error(), http.StatusInternalServerError)
		t.code = http.StatusInternalServerError
		return hijackError
	}

	go t.transfer(clientConn, t.DestConn)
	go t.transfer(t.DestConn, clientConn)

	return nil
}

func (t *responseWriterTunnel) transfer(src io.ReadCloser, dst io.WriteCloser) {
	var nr, nw int
	var err error
	buf := make([]byte, bufferSize)

	for {
		ctx, cancel := contextWithTimeout(t.ReadTimeout)
		go func() {
			nr, err = src.Read(buf)
			cancel()
		}()

		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded || err != nil {
			break
		}

		if nr > 0 {
			ctx, cancel := contextWithTimeout(t.WriteTimeout)
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

	closeFile(dst)
	closeFile(src)
}

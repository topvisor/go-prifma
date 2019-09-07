package proxy

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const bufferSize = 1024 * 32

type responseWriterTunnel struct {
	Dialer       *dialer
	Proxy        *Proxy
	Request      *http.Request
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	lAddr net.Addr
	rAddr net.Addr
	code  int
}

func (t *responseWriterTunnel) GetCode() int {
	return t.code
}

func (t *responseWriterTunnel) GetLAddr() net.Addr {
	return t.lAddr
}

func (t *responseWriterTunnel) GetRAddr() net.Addr {
	return t.rAddr
}

func (t *responseWriterTunnel) Write(rw http.ResponseWriter) error {
	var destConn net.Conn
	var err error

	if t.Proxy != nil {
		if destConn, err = t.Proxy.connect(t.Request, t.Dialer); err != nil {
			http.Error(rw, err.Error(), http.StatusBadGateway)
			t.code = http.StatusBadGateway

			return err
		}
	}
	if destConn == nil {
		destUrl := &url.URL{Host: t.Request.Host}
		if destUrl.Port() == "" {
			destUrl.Host = net.JoinHostPort(t.Request.Host, "443")
		}
		if destConn, err = t.Dialer.connect(destUrl); err != nil {
			http.Error(rw, err.Error(), http.StatusBadGateway)
			t.code = http.StatusBadGateway

			return err
		}
	}

	rw.WriteHeader(http.StatusOK)
	t.code = http.StatusOK

	clientConn, _, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		closeFile(destConn)

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		t.code = http.StatusInternalServerError

		return err
	}

	go t.transfer(clientConn, destConn)
	go t.transfer(destConn, clientConn)

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

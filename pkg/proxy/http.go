package proxy

import (
	"io"
	"net"
	"net/url"
	"time"
)

func connectToHost(network string, host string, laddr net.Addr, dialTimeout time.Duration) (net.Conn, error) {
	dialer := net.Dialer{
		Timeout:   dialTimeout,
		LocalAddr: laddr,
	}

	return dialer.Dial(network, host)
}

func connectToUrl(network string, url url.URL, laddr net.Addr, dialTimeout time.Duration) (net.Conn, error) {
	host := url.Host
	if url.Port() == "" {
		host = net.JoinHostPort(host, "80")
	}

	return connectToHost(network, host, laddr, dialTimeout)
}

func transfer(src io.ReadCloser, dst io.WriteCloser) {
	_, _ = io.Copy(dst, src)
	_ = src.Close()
	_ = dst.Close()
}

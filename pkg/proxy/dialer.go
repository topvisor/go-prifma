package proxy

import (
	"fmt"
	"net"
	"net/url"
)

type dialer struct {
	lIpV4 net.IP
	lIpV6 net.IP

	net.Dialer
}

func (t *dialer) connect(url *url.URL) (net.Conn, error) {
	host := url.Hostname()
	if host == "" {
		return nil, fmt.Errorf("unavailable url: \"%s\"", url.String())
	}

	port := url.Port()
	if port == "" {
		switch url.Scheme {
		case "http":
			fallthrough
		case "ws":
			port = "80"
		case "https":
			fallthrough
		case "wss":
			port = "443"
		default:
			return nil, fmt.Errorf("unavailable url: \"%s\"", url.String())
		}
	}

	laddr, err := t.selectLAddr(host)
	if err != nil {
		return nil, err
	}

	t.LocalAddr = laddr

	return t.Dial("tcp", net.JoinHostPort(host, port))
}

func (t *dialer) selectLAddr(host string) (net.Addr, error) {
	if t.lIpV4 == nil && t.lIpV6 == nil {
		return nil, nil
	}
	if t.lIpV4 != nil && t.lIpV6 == nil {
		return &net.TCPAddr{IP: t.lIpV4}, nil
	}
	if t.lIpV6 != nil && t.lIpV4 == nil {
		return &net.TCPAddr{IP: t.lIpV6}, nil
	}

	destIps, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	if len(destIps) > 0 {
		if destIps[0].To4() != nil {
			return &net.TCPAddr{IP: t.lIpV4}, nil
		} else {
			return &net.TCPAddr{IP: t.lIpV6}, nil
		}
	} else {
		return nil, fmt.Errorf("unreachable host: \"%s\"", host)
	}
}

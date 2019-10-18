package prifma

import (
	"fmt"
	"net"
	"net/url"
)

// dialer is a net.Dialer which select a suitable outgoing ip address depending on a request's destination domain
// lIpV4 and lIpV6 set a ip addresses from which to select
type dialer struct {
	lIpV4 net.IP
	lIpV6 net.IP
}

// connect connects to the url
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
		case "tunnel":
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

	dialer := &net.Dialer{
		LocalAddr: laddr,
	}

	return dialer.Dial("tcp", net.JoinHostPort(host, port))
}

// selectLAddr select a suitable outgoing ip address depending on a request's destination domain
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

// ipsString generates a string which represents ip addresses which passed to dialer
func (t *dialer) ipsString() string {
	ipV4Str := ""
	if t.lIpV4 != nil {
		ipV4Str = t.lIpV4.String()
	}

	ipV6Str := ""
	if t.lIpV6 != nil {
		ipV6Str = t.lIpV6.String()
	}

	return fmt.Sprintf("%s:%s", ipV4Str, ipV6Str)
}

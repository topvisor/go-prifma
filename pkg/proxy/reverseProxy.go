package proxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type reverseProxy struct {
	dialer *dialer

	httputil.ReverseProxy
}

func newReverseProxy(handler *Handler, dialer *dialer) *reverseProxy {
	reverseProxy := &reverseProxy{
		dialer: dialer,
	}

	reverseProxy.Director = removeProxyHeaders
	reverseProxy.Transport = &http.Transport{
		Proxy:              http.ProxyURL(handler.Proxy.Url),
		ProxyConnectHeader: handler.Proxy.ProxyHeaders,
		DialContext:        reverseProxy.DialContext,
	}
	reverseProxy.FlushInterval = -1
	reverseProxy.ErrorLog = handler.ErrorLogger.logger

	return reverseProxy
}

func (t *reverseProxy) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	destUrl := &url.URL{Host: addr}
	if destUrl.Port() == "" {
		destUrl.Host = net.JoinHostPort(addr, "80")
	}

	return t.dialer.connect(destUrl)
}

var proxyHeadersRegexp, _ = regexp.Compile("^(?i)proxy-")

func removeProxyHeaders(req *http.Request) {
	for key := range req.Header {
		if proxyHeadersRegexp.MatchString(key) {
			req.Header.Del(key)
		}
	}
}

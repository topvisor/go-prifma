package prifma_new

import (
	"context"
	"net"
	"net/http"
	"net/url"
)

type ProxyFunc func(*http.Request) (*url.URL, error)

type HandleRequestResult interface {
	SetRequest(req *http.Request)
	SetDialer(dialer Dialer)
	SetProxy(proxy ProxyFunc)
	SetProxyConnectHeader(header http.Header)
	SetResponse(resp Response)

	GetRequest() *http.Request
	GetDialer() Dialer
	GetProxy() ProxyFunc
	GetProxyConnectHeader() http.Header
	GetResponse() Response
	GetServer() Server

	GetRoundTripper() http.RoundTripper
}

func NewHandleRequestResult(req *http.Request, server Server) HandleRequestResult {
	t := &DefaultHandleRequestResult{
		Server:  server,
		Request: req,
		Dialer:  NewDialer(),
	}

	t.Transport = &http.Transport{
		DialContext:           t.DialContext,
		IdleConnTimeout:       server.GetIdleTimeout(),
		ResponseHeaderTimeout: server.GetWriteTimeout(),
	}

	return t
}

type DefaultHandleRequestResult struct {
	Server    Server
	Request   *http.Request
	Response  Response
	Dialer    Dialer
	Transport *http.Transport
}

func (t *DefaultHandleRequestResult) SetRequest(req *http.Request) {
	t.Request = req
}

func (t *DefaultHandleRequestResult) SetDialer(dialer Dialer) {
	t.Dialer = dialer
}

func (t *DefaultHandleRequestResult) SetProxy(proxy ProxyFunc) {
	t.Transport.Proxy = proxy
}

func (t *DefaultHandleRequestResult) SetProxyConnectHeader(header http.Header) {
	t.Transport.ProxyConnectHeader = header
}

func (t *DefaultHandleRequestResult) SetResponse(resp Response) {
	t.Response = resp
}

func (t *DefaultHandleRequestResult) GetRequest() *http.Request {
	return t.Request
}

func (t *DefaultHandleRequestResult) GetDialer() Dialer {
	return t.Dialer
}

func (t *DefaultHandleRequestResult) GetProxy() ProxyFunc {
	return t.Transport.Proxy
}

func (t *DefaultHandleRequestResult) GetProxyConnectHeader() http.Header {
	return t.Transport.ProxyConnectHeader
}

func (t *DefaultHandleRequestResult) GetResponse() Response {
	return t.Response
}

func (t *DefaultHandleRequestResult) GetServer() Server {
	return t.Server
}

func (t *DefaultHandleRequestResult) GetRoundTripper() http.RoundTripper {
	return t.Transport
}

func (t *DefaultHandleRequestResult) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return t.Dialer.DialContext(ctx, network, addr)
}

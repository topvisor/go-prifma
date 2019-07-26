package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type reverseProxyRequestData struct {
	Response *http.Response
	Error    error
}

type reverseProxy struct {
	RequestData map[uint64]*reverseProxyRequestData

	dialer *dialer

	httputil.ReverseProxy
}

func newReverseProxy(handler *Handler, dialer *dialer) *reverseProxy {
	reverseProxy := &reverseProxy{
		RequestData: make(map[uint64]*reverseProxyRequestData),
		dialer:      dialer,
	}

	transport := &http.Transport{
		DialContext: reverseProxy.dialContext,
	}
	if handler.Proxy != nil {
		transport.Proxy = http.ProxyURL(handler.Proxy.Url)
		transport.ProxyConnectHeader = handler.Proxy.ProxyHeaders
	}

	reverseProxy.Director = removeProxyHeaders
	reverseProxy.ModifyResponse = reverseProxy.saveResponse
	reverseProxy.ErrorHandler = reverseProxy.errorHandler
	reverseProxy.Transport = transport
	reverseProxy.FlushInterval = -1

	return reverseProxy
}

func (t *reverseProxy) dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	destUrl := &url.URL{Host: addr}
	if destUrl.Port() == "" {
		destUrl.Host = net.JoinHostPort(addr, "80")
	}

	return t.dialer.connect(destUrl)
}

func (t *reverseProxy) saveResponse(resp *http.Response) error {
	reqIdInterface := resp.Request.Context().Value(keyReqId)
	if reqIdInterface == nil {
		return nil
	}

	reqId := reqIdInterface.(uint64)
	if _, exists := t.RequestData[reqId]; exists {
		t.RequestData[reqId].Response = resp
	} else {
		t.RequestData[reqId] = &reverseProxyRequestData{Response: resp}
	}

	return nil
}

func (t *reverseProxy) errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	reqIdInterface := req.Context().Value(keyReqId)
	if reqIdInterface == nil {
		return
	}

	reqId := reqIdInterface.(uint64)
	if _, exists := t.RequestData[reqId]; !exists {
		t.RequestData[reqId] = &reverseProxyRequestData{Response: new(http.Response)}
	}

	switch err {
	case context.DeadlineExceeded:
		http.Error(rw, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		t.RequestData[reqId].Response.StatusCode = http.StatusGatewayTimeout
		t.RequestData[reqId].Error = fmt.Errorf("%d, %s", http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))
	case context.Canceled:
		http.Error(rw, StatusTextClientClosedRequest, StatusClientClosedRequest)
		t.RequestData[reqId].Response.StatusCode = StatusClientClosedRequest
		t.RequestData[reqId].Error = fmt.Errorf("%d, %s", StatusClientClosedRequest, StatusTextClientClosedRequest)
		// ###
	}
}

var proxyHeadersRegexp, _ = regexp.Compile("^(?i)proxy-")

func removeProxyHeaders(req *http.Request) {
	for key := range req.Header {
		if proxyHeadersRegexp.MatchString(key) {
			req.Header.Del(key)
		}
	}
}

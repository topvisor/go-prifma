package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/topvisor/prifma/pkg/prifma"
	"github.com/topvisor/prifma/pkg/utils"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
)

type ResponseReverseProxy struct {
	RoundTrippers RoundTrippersMap
	ResponseCode  int
	Error         error
	LAddr         net.Addr
	RAddr         net.Addr
}

func NewResponseReverseProxy(roundTrippers RoundTrippersMap) prifma.Response {
	return &ResponseReverseProxy{
		ResponseCode:  http.StatusInternalServerError,
		Error:         errors.New(http.StatusText(http.StatusInternalServerError)),
		RoundTrippers: roundTrippers,
	}
}

func (t *ResponseReverseProxy) Write(rw http.ResponseWriter, result prifma.HandleRequestResult) error {
	reverseProxy := &httputil.ReverseProxy{
		Director:       utils.RemoveProxyHeaders,
		Transport:      t.RoundTrippers.Get(result),
		FlushInterval:  -1,
		ModifyResponse: t.SaveResponse,
		ErrorHandler:   t.ErrorHandler,
	}

	req := result.GetRequest().WithContext(
		httptrace.WithClientTrace(
			result.GetRequest().Context(),
			&httptrace.ClientTrace{
				GotConn: func(info httptrace.GotConnInfo) {
					t.LAddr = info.Conn.LocalAddr()
					t.RAddr = info.Conn.RemoteAddr()
				},
			},
		),
	)

	reverseProxy.ServeHTTP(rw, req)

	return t.Error
}

func (t *ResponseReverseProxy) GetCode() int {
	return t.ResponseCode
}

func (t *ResponseReverseProxy) GetLAddr() net.Addr {
	return t.LAddr
}

func (t *ResponseReverseProxy) GetRAddr() net.Addr {
	return t.RAddr
}

func (t *ResponseReverseProxy) SaveResponse(resp *http.Response) error {
	t.ResponseCode = resp.StatusCode
	t.Error = nil

	return nil
}

func (t *ResponseReverseProxy) ErrorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	t.Error = nil

	switch err {
	case context.DeadlineExceeded:
		http.Error(rw, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
		t.ResponseCode = http.StatusGatewayTimeout
		t.Error = fmt.Errorf("%d, %s", http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))
	case context.Canceled:
		http.Error(rw, prifma.StatusTextClientClosedRequest, prifma.StatusClientClosedRequest)
		t.ResponseCode = prifma.StatusClientClosedRequest
	default:
		switch err := err.(type) {
		case *net.OpError:
			if err.Op == "dial" {
				http.Error(rw, err.Error(), http.StatusBadGateway)
				t.ResponseCode = http.StatusBadGateway
			} else {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				t.ResponseCode = http.StatusInternalServerError
				t.Error = fmt.Errorf("%d, %s", http.StatusInternalServerError, err.Error())
			}
		}
	}
}

package proxy

import (
	"context"
	"errors"
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"io"
	"net"
	"net/http"
	"time"
)

type responseContext struct {
	writer responseWriter
	ctx    context.Context
	cancel context.CancelFunc
}

func newResponseTimeout(timeout time.Duration) *responseContext {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)

	return &responseContext{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *responseContext) Cancel() {
	t.cancel()
}

type responseWriter interface {
	Write(rw http.ResponseWriter) error
	GetCode() int
}

type responseWriterError struct {
	Code  int
	Error string
}

func (t *responseWriterError) GetCode() int {
	return t.Code
}

func (t *responseWriterError) Write(rw http.ResponseWriter) error {
	errorString := t.Error
	if errorString == "" {
		errorString = http.StatusText(t.Code)
	}

	http.Error(rw, errorString, t.Code)

	return fmt.Errorf("%d %s", t.Code, t.Error)
}

type responseWriteRequireAuth struct {
	Request   *http.Request
	BasicAuth *auth.BasicAuth
}

func (t *responseWriteRequireAuth) GetCode() int {
	return t.BasicAuth.Headers.UnauthCode
}

func (t *responseWriteRequireAuth) Write(rw http.ResponseWriter) error {
	t.BasicAuth.RequireAuth(rw, t.Request)

	return nil
}

type responseWriterTunnel struct {
	DestConn net.Conn

	code int
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

	go transfer(clientConn, t.DestConn)
	go transfer(t.DestConn, clientConn)

	return nil
}

type responseWriteReverseProxy struct {
	Request      *http.Request
	ReverseProxy *reverseProxy

	code int
}

func (t *responseWriteReverseProxy) GetCode() int {
	return t.code
}

func (t *responseWriteReverseProxy) Write(rw http.ResponseWriter) error {
	t.ReverseProxy.ServeHTTP(rw, t.Request)

	reqIdInterface := t.Request.Context().Value(keyReqId)
	if reqIdInterface == nil {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		t.code = http.StatusInternalServerError
		return errors.New("request id must not be <nil>")
	}

	reqId := reqIdInterface.(uint64)
	reqData, reqDataExists := t.ReverseProxy.RequestData[reqId]
	if !reqDataExists {
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		t.code = http.StatusInternalServerError
		return errors.New("request data not exists")
	}

	t.code = reqData.Response.StatusCode

	return nil
}

func transfer(src io.ReadCloser, dst io.WriteCloser) {
	_, _ = io.Copy(dst, src)
	_ = src.Close()
	_ = dst.Close()
}

package proxy

import (
	"errors"
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"net/http"
)

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

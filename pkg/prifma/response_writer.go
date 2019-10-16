package prifma

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"net"
	"net/http"
)

type responseWriter interface {
	Write(rw http.ResponseWriter) error
	GetCode() int
	GetLAddr() net.Addr
	GetRAddr() net.Addr
}

type responseWriterError struct {
	Code  int
	Error string
}

func (t *responseWriterError) GetCode() int {
	return t.Code
}

func (t *responseWriterError) GetLAddr() net.Addr {
	return nil
}

func (t *responseWriterError) GetRAddr() net.Addr {
	return nil
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

func (t *responseWriteRequireAuth) GetLAddr() net.Addr {
	return nil
}

func (t *responseWriteRequireAuth) GetRAddr() net.Addr {
	return nil
}

func (t *responseWriteRequireAuth) Write(rw http.ResponseWriter) error {
	t.BasicAuth.RequireAuth(rw, t.Request)

	return nil
}

package prifma_new

import (
	"fmt"
	"net"
	"net/http"
)

type Response interface {
	Write(rw http.ResponseWriter) error
	GetCode() int
	GetLAddr() net.Addr
	GetRAddr() net.Addr
}

type ResponseError struct {
	Code  int
	Error string
}

func NewResponseError(code int, error string) Response {
	return &ResponseError{
		Code:  code,
		Error: error,
	}
}

func (t *ResponseError) Write(rw http.ResponseWriter) error {
	errStr := t.Error
	if errStr == "" {
		errStr = http.StatusText(t.Code)
	}

	http.Error(rw, errStr, t.Code)

	return fmt.Errorf("%d %s", t.Code, errStr)
}

func (t *ResponseError) GetCode() int {
	return t.Code
}

func (*ResponseError) GetLAddr() net.Addr {
	return nil
}

func (*ResponseError) GetRAddr() net.Addr {
	return nil
}

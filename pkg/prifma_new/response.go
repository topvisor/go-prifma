package prifma_new

import (
	"net"
	"net/http"
)

type Response interface {
	Write(rw http.ResponseWriter) error
	GetCode() int
	GetLAddr() net.Addr
	GetRAddr() net.Addr
}

type TestResponse struct {
}

func (t *TestResponse) Write(rw http.ResponseWriter) error {
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write(([]byte)(http.StatusText(http.StatusOK)))

	return err
}

func (t *TestResponse) GetCode() int {
	return http.StatusOK
}

func (t *TestResponse) GetLAddr() net.Addr {
	return nil
}

func (t *TestResponse) GetRAddr() net.Addr {
	return nil
}

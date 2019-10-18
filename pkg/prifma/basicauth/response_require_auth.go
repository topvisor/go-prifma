package basicauth

import (
	auth "github.com/abbot/go-http-auth"
	"github.com/topvisor/prifma/pkg/prifma"
	"net"
	"net/http"
)

type ResponseRequireAuth struct {
	Request *http.Request
}

func NewResponseRequireAuth(req *http.Request) prifma.Response {
	return &ResponseRequireAuth{
		Request: req,
	}
}

func (t *ResponseRequireAuth) Write(rw http.ResponseWriter, _ prifma.HandleRequestResult) error {
	basicAuth := auth.BasicAuth{
		Headers: auth.ProxyHeaders,
	}

	basicAuth.RequireAuth(rw, t.Request)

	return nil
}

func (t *ResponseRequireAuth) GetCode() int {
	return auth.ProxyHeaders.UnauthCode
}

func (t *ResponseRequireAuth) GetLAddr() net.Addr {
	return nil
}

func (t *ResponseRequireAuth) GetRAddr() net.Addr {
	return nil
}

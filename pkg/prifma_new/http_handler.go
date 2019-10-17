package prifma_new

import (
	"context"
	"net"
	"net/http"
)

const (
	StatusClientClosedRequest     = 499
	StatusTextClientClosedRequest = "Client Closed Request"
)

const (
	CtxKeyTransport = iota
)

type HttpHandler struct {
	Server Server
}

func NewHttpHandler(server Server) http.Handler {
	return &HttpHandler{
		Server: server,
	}
}

func (t *HttpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	modules := t.Server.GetModulesManager().GetModulesForRequest(req)

	for _, module := range modules {
		if handler, ok := module.(BeforeHandleRequestModule); ok {
			if err := handler.BeforeHandleRequest(req); err != nil {
				t.Server.GetErrorLog().Println(err)
			}
		}
	}

	result := NewHandleRequestResult(req)

	for _, module := range modules {
		if handler, ok := module.(HandleRequestModule); ok {
			var err error
			if result, err = handler.HandleRequest(result); err != nil {
				t.Server.GetErrorLog().Println(err)
			}
			if result.GetResponse() != nil {
				break
			}
		}
	}

	if req.Context().Err() != nil {
		resp = NewResponseError(StatusClientClosedRequest, StatusTextClientClosedRequest)
	}
	if resp == nil {
		resp = NewResponseError(http.StatusInternalServerError, "")
	}
	if err := resp.Write(rw); err != nil {
		t.Server.GetErrorLog().Println(err)
	}

	for _, module := range modules {
		if handler, ok := module.(AfterWriteResponseModule); ok {
			if err := handler.AfterWriteResponse(req, resp); err != nil {
				t.Server.GetErrorLog().Println(err)
			}
		}
	}
}

func (t *HttpHandler) SetTransport(req *http.Request) *http.Request {
	transport := &http.Transport{
		DialContext: net.Dialer{
			Timeout: t.Server.GetWriteTimeout(),
		}.DialContext,
	}
}

func (t *HttpHandler) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {

}

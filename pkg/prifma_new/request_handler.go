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

type RequestHandler struct {
	Server Server
}

func NewRequestHandler(server Server) http.Handler {
	return &RequestHandler{
		Server: server,
	}
}

func (t *RequestHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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

func (t *RequestHandler) SetTransport(req *http.Request) *http.Request {
	transport := &http.Transport{
		DialContext: net.Dialer{
			Timeout: t.Server.GetWriteTimeout(),
		}.DialContext,
	}
}

func (t *RequestHandler) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {

}

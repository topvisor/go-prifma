package prifma

import (
	"net/http"
)

const (
	StatusClientClosedRequest     = 499
	StatusTextClientClosedRequest = "Client Closed Request"
)

type RequestHandler struct {
	Server Server
}

func NewRequestHandler(server Server) *RequestHandler {
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

	var result HandleRequestResult = NewHandleRequestResult(req, t.Server)

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
		result.SetResponse(NewResponseError(StatusClientClosedRequest, StatusTextClientClosedRequest))
	}
	if result.GetResponse() == nil {
		result.SetResponse(NewResponseError(http.StatusInternalServerError, ""))
	}
	if err := result.GetResponse().Write(rw, result); err != nil {
		t.Server.GetErrorLog().Println(err)
	}

	for _, module := range modules {
		if handler, ok := module.(AfterWriteResponseModule); ok {
			if err := handler.AfterWriteResponse(req, result.GetResponse()); err != nil {
				t.Server.GetErrorLog().Println(err)
			}
		}
	}
}

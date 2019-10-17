package prifma_new

import "net/http"

const (
	StatusClientClosedRequest     = 499
	StatusTextClientClosedRequest = "Client Closed Request"
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

	var resp Response

	for _, module := range modules {
		if handler, ok := module.(HandleRequestModule); ok {
			var err error
			if req, resp, err = handler.HandleRequest(req); err != nil {
				t.Server.GetErrorLog().Println(err)
			}
			if resp != nil {
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

package prifma_new

import "net/http"

type HttpHandler struct {
	Server Server
}

func NewHttpHandler(server Server) http.Handler {
	return &HttpHandler{
		Server: server,
	}
}

func (t *HttpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, module := range t.Server.GetModules() {
		if handler, ok := module.(BeforeHandleRequestHandler); ok {
			_ = handler.BeforeHandleRequest(req)
		}
	}

	resp := new(TestResponse)
	_ = resp.Write(rw)

	for _, module := range t.Server.GetModules() {
		if handler, ok := module.(AfterWriteResponseHandler); ok {
			_ = handler.AfterWriteResponse(req, resp)
		}
	}
}

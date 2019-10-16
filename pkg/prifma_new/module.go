package prifma_new

import (
	"github.com/topvisor/prifma/pkg/conf"
	"net/http"
)

type Module interface {
	conf.Block
}

type BeforeHandleRequestHandler interface {
	BeforeHandleRequest(req *http.Request) error
}

type AfterWriteResponseHandler interface {
	AfterWriteResponse(req *http.Request, resp Response) error
}

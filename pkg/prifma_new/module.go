package prifma_new

import (
	"github.com/topvisor/prifma/pkg/conf"
	"net/http"
)

type Module interface {
	GetDirective() string
	Clone() Module

	conf.Block
}

type BeforeHandleRequestModule interface {
	BeforeHandleRequest(req *http.Request) error
}

type HandleRequestModule interface {
	HandleRequest(result HandleRequestResult) (HandleRequestResult, error)
}

type AfterWriteResponseModule interface {
	AfterWriteResponse(req *http.Request, resp Response) error
}

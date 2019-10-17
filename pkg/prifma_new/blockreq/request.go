package blockreq

import (
	"github.com/topvisor/prifma/pkg/prifma_new"
	"net"
	"net/http"
)

type ResponseLocked struct{}

func NewResponseLocked() prifma_new.Response {
	return new(ResponseLocked)
}

func (t *ResponseLocked) Write(rw http.ResponseWriter, _ prifma_new.HandleRequestResult) error {
	http.Error(rw, http.StatusText(t.GetCode()), t.GetCode())

	return nil
}

func (t *ResponseLocked) GetCode() int {
	return http.StatusLocked
}

func (t *ResponseLocked) GetLAddr() net.Addr {
	return nil
}

func (t *ResponseLocked) GetRAddr() net.Addr {
	return nil
}

package blockreq

import (
	"github.com/topvisor/go-prifma/pkg/prifma"
	"net"
	"net/http"
)

type ResponseLocked struct{}

func NewResponseLocked() *ResponseLocked {
	return new(ResponseLocked)
}

func (t *ResponseLocked) Write(rw http.ResponseWriter, _ prifma.HandleRequestResult) error {
	rw.Header().Add("X-Prifma-Error", http.StatusText(t.GetCode()))
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

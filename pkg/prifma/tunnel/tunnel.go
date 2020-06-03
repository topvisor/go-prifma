package tunnel

import (
	"github.com/topvisor/go-prifma/pkg/conf"
	"github.com/topvisor/go-prifma/pkg/prifma"
	"net/http"
)

const ModuleDirective = "tunnel"

type Tunnel struct {
}

func New() *Tunnel {
	return &Tunnel{}
}

func (t *Tunnel) HandleRequest(result prifma.HandleRequestResult) (prifma.HandleRequestResult, error) {
	if result.GetRequest().Method == http.MethodConnect {
		result.SetResponse(NewResponseTunnel())
	}

	return result, nil
}

func (t *Tunnel) GetDirective() string {
	return ModuleDirective
}

func (t *Tunnel) Clone() prifma.Module {
	return t
}

func (t *Tunnel) Call(command conf.Command) error {
	return conf.NewErrCommandName(command)
}

func (t *Tunnel) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommandName(command)
}

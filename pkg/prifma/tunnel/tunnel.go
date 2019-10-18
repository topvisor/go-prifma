package tunnel

import (
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma"
	"net/http"
)

const ModuleDirective = "tunnel"

type Tunnel struct {
}

func New() prifma.Module {
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
	return conf.NewErrCommand(command)
}

func (t *Tunnel) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}
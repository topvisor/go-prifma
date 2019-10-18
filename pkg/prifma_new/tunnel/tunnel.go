package tunnel

import (
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma_new"
	"net/http"
)

const ModuleDirective = "tunnel"

type Tunnel struct {
}

func New() prifma_new.Module {
	return &Tunnel{}
}

func (t *Tunnel) HandleRequest(result prifma_new.HandleRequestResult) (prifma_new.HandleRequestResult, error) {
	if result.GetRequest().Method == http.MethodConnect {
		result.SetResponse(NewResponseTunnel())
	}

	return result, nil
}

func (t *Tunnel) GetDirective() string {
	return ModuleDirective
}

func (t *Tunnel) Clone() prifma_new.Module {
	return t
}

func (t *Tunnel) Call(command conf.Command) error {
	return conf.NewErrCommand(command)
}

func (t *Tunnel) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}

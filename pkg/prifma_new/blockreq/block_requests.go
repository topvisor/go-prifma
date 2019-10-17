package blockreq

import (
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma_new"
	"net/http"
)

const ModuleDirective = "block_requests"

type BlockRequests struct {
	Enabled bool
}

func New() prifma_new.Module {
	return new(BlockRequests)
}

func (t *BlockRequests) HandleRequest(req *http.Request) (*http.Request, prifma_new.Response, error) {
	if t.Enabled {
		return req, NewResponseLocked(), nil
	} else {
		return req, nil, nil
	}
}

func (t *BlockRequests) Off() error {
	t.Enabled = true

	return nil
}

func (t *BlockRequests) On() error {
	t.Enabled = false

	return nil
}

func (t *BlockRequests) GetDirective() string {
	return ModuleDirective
}

func (t *BlockRequests) Clone() prifma_new.Module {
	clone := *t

	return &clone
}

func (t *BlockRequests) Call(command conf.Command) error {
	if command.GetName() != ModuleDirective || len(command.GetArgs()) != 1 {
		return conf.NewErrCommand(command)
	}

	switch command.GetArgs()[0] {
	case "off":
		return t.Off()
	case "on":
		return t.On()
	}

	return conf.NewErrCommand(command)
}

func (t *BlockRequests) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}

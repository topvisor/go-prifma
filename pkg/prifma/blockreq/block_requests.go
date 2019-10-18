package blockreq

import (
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma"
)

const ModuleDirective = "block_requests"

type BlockRequests struct {
	Enabled bool
}

func New() prifma.Module {
	return new(BlockRequests)
}

func (t *BlockRequests) HandleRequest(result prifma.HandleRequestResult) (prifma.HandleRequestResult, error) {
	if t.Enabled {
		result.SetResponse(NewResponseLocked())
	}

	return result, nil
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

func (t *BlockRequests) Clone() prifma.Module {
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

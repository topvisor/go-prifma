package blockreq

import (
	"github.com/topvisor/go-prifma/pkg/conf"
	"github.com/topvisor/go-prifma/pkg/prifma"
)

const ModuleDirective = "block_requests"

type BlockRequests struct {
	Enabled bool
}

func New() *BlockRequests {
	return new(BlockRequests)
}

func (t *BlockRequests) HandleRequest(result prifma.HandleRequestResult) (prifma.HandleRequestResult, error) {
	if t.Enabled {
		result.SetResponse(NewResponseLocked())
	}

	return result, nil
}

func (t *BlockRequests) Off() error {
	t.Enabled = false

	return nil
}

func (t *BlockRequests) On() error {
	t.Enabled = true

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
	if command.GetName() != ModuleDirective {
		return conf.NewErrCommandName(command)
	}

	if len(command.GetArgs()) != 1 {
		return conf.NewErrCommandArgsNumber(command)
	}

	arg := command.GetArgs()[0]

	switch arg {
	case "off":
		return t.Off()
	case "on":
		return t.On()
	}

	return conf.NewErrCommandArg(command, arg)
}

func (t *BlockRequests) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommandMustHaveNoBlock(command)
}

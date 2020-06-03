package outgoingip

import (
	"github.com/topvisor/go-prifma/pkg/conf"
)

type IpArray interface {
	AddIp(ip string) error
}

type ConfBlock struct {
	IpArray IpArray
}

func NewConfBlock(ipArray IpArray) *ConfBlock {
	return &ConfBlock{
		IpArray: ipArray,
	}
}

func (t *ConfBlock) Call(command conf.Command) (err error) {
	if len(command.GetArgs()) != 0 {
		return conf.NewErrCommandArgsNumber(command)
	}

	if err = t.IpArray.AddIp(command.GetName()); err != nil {
		err = conf.NewErrCommand(command, err.Error())
	}

	return err
}

func (t *ConfBlock) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommandMustHaveNoBlock(command)
}

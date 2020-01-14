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

func NewConfBlock(ipArray IpArray) conf.Block {
	return &ConfBlock{
		IpArray: ipArray,
	}
}

func (t *ConfBlock) Call(command conf.Command) error {
	return t.IpArray.AddIp(command.GetName())
}

func (t *ConfBlock) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}

package outgoingip

import (
	"github.com/topvisor/prifma/pkg/conf"
)

type IpArray interface {
	AddIp(ip string) error
}

type OutgoingIpBlock struct {
	IpArray IpArray
}

func NewOutgoingIpBlock(outgoingIp IpArray) conf.Block {
	return &OutgoingIpBlock{
		IpArray: outgoingIp,
	}
}

func (t *OutgoingIpBlock) Call(command conf.Command) error {
	return t.IpArray.AddIp(command.GetName())
}

func (t *OutgoingIpBlock) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}

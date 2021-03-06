package outgoingip

import (
	"fmt"
	"github.com/topvisor/go-prifma/pkg/conf"
	"github.com/topvisor/go-prifma/pkg/prifma"
	"math/rand"
	"net"
)

const ModuleDirective = "outgoing_ip"

type OutgoingIp struct {
	IpsV4 []net.IP
	IpsV6 []net.IP
}

func New() *OutgoingIp {
	return &OutgoingIp{
		IpsV4: make([]net.IP, 0),
		IpsV6: make([]net.IP, 0),
	}
}

func (t *OutgoingIp) HandleRequest(result prifma.HandleRequestResult) (prifma.HandleRequestResult, error) {
	ipsV4Len := len(t.IpsV4)
	ipsV6Len := len(t.IpsV6)

	if ipsV4Len == 0 && ipsV6Len == 0 {
		return result, nil
	}

	result.GetDialer().SetIpV4(nil)
	result.GetDialer().SetIpV6(nil)

	if ipsV4Len != 0 {
		result.GetDialer().SetIpV4(t.IpsV4[rand.Intn(ipsV4Len)])
	}
	if ipsV6Len != 0 {
		result.GetDialer().SetIpV6(t.IpsV6[rand.Intn(ipsV6Len)])
	}

	return result, nil
}

func (t *OutgoingIp) Off() error {
	t.IpsV4 = make([]net.IP, 0)
	t.IpsV6 = make([]net.IP, 0)

	return nil
}

func (t *OutgoingIp) SetIps(ips []string) error {
	t.IpsV4 = make([]net.IP, 0, len(ips))
	t.IpsV6 = make([]net.IP, 0, len(ips))

	for _, ip := range ips {
		if err := t.AddIp(ip); err != nil {
			return err
		}
	}

	return nil
}

func (t *OutgoingIp) AddIp(ipStr string) error {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("wrong outgoing ip - '%s'", ipStr)
	}

	if ipV4 := ip.To4(); ipV4 != nil {
		t.IpsV4 = append(t.IpsV4, ipV4)
	} else {
		t.IpsV6 = append(t.IpsV6, ip)
	}

	return nil
}

func (t *OutgoingIp) GetDirective() string {
	return ModuleDirective
}

func (t *OutgoingIp) Clone() prifma.Module {
	clone := *t

	return &clone
}

func (t *OutgoingIp) Call(command conf.Command) (err error) {
	if command.GetName() != ModuleDirective {
		return conf.NewErrCommandName(command)
	}

	args := command.GetArgs()
	if len(args) == 0 {
		return conf.NewErrCommandArgsNumber(command)
	}

	if len(args) == 1 && args[0] == "off" {
		return t.Off()
	}

	if err = t.SetIps(args); err != nil {
		err = conf.NewErrCommand(command, err.Error())
	}

	return err
}

func (t *OutgoingIp) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() != ModuleDirective {
		return nil, conf.NewErrCommandName(command)
	}

	if len(command.GetArgs()) != 0 {
		return nil, conf.NewErrCommandArgsNumber(command)
	}

	return NewConfBlock(t), nil
}

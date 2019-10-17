package outgoingip

import (
	"context"
	"fmt"
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma_new"
	"math/rand"
	"net"
	"net/http"
)

const ModuleDirective = "outgoing_ip"

type OutgoingIp struct {
	IpsV4 []net.IP
	IpsV6 []net.IP
}

func New() prifma_new.Module {
	return &OutgoingIp{
		IpsV4: make([]net.IP, 0),
		IpsV6: make([]net.IP, 0),
	}
}

func (t *OutgoingIp) HandleRequest(req *http.Request) (*http.Request, prifma_new.Response, error) {
	ipsV4Len := len(t.IpsV4)
	ipsV6Len := len(t.IpsV6)

	if ipsV4Len == 0 && ipsV6Len == 0 {
		return req, nil, nil
	}

	ctx := req.Context()

	if ipsV4Len == 0 {
		ctx = context.WithValue(ctx, prifma_new.CtxOutgoingIpV4, nil)
	} else {
		ctx = context.WithValue(ctx, prifma_new.CtxOutgoingIpV4, t.IpsV4[rand.Intn(ipsV4Len)])
	}
	if ipsV6Len == 0 {
		ctx = context.WithValue(ctx, prifma_new.CtxOutgoingIpV6, nil)
	} else {
		ctx = context.WithValue(ctx, prifma_new.CtxOutgoingIpV6, t.IpsV6[rand.Intn(ipsV6Len)])
	}

	return req.WithContext(ctx), nil, nil
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
		return fmt.Errorf("wrong outgoing ip: '%s'", ipStr)
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

func (t *OutgoingIp) Clone() prifma_new.Module {
	clone := *t

	return &clone
}

func (t *OutgoingIp) Call(command conf.Command) error {
	args := command.GetArgs()
	if command.GetName() != ModuleDirective || len(args) == 0 {
		return conf.NewErrCommand(command)
	}
	if len(args) == 1 && args[1] == "off" {
		return t.Off()
	}

	return t.SetIps(args)
}

func (t *OutgoingIp) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() != ModuleDirective || len(command.GetArgs()) != 0 {
		return nil, conf.NewErrCommand(command)
	}

	return NewOutgoingIpBlock(t), nil
}

package useipheader

import (
	"context"
	"fmt"
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma_new"
	"net"
	"net/http"
)

const (
	ModuleDirective = "use_ip_header"
	HeaderName      = "Proxy-Use-Ip"
)

type UseIpHeader struct {
	Enabled bool
}

func New() prifma_new.Module {
	return new(UseIpHeader)
}

func (t *UseIpHeader) HandleRequest(req *http.Request) (*http.Request, prifma_new.Response, error) {
	if !t.Enabled {
		return req, nil, nil
	}

	ipStr := req.Header.Get(HeaderName)
	if ipStr == "" {
		return req, nil, nil
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return req, prifma_new.NewResponseError(http.StatusBadRequest, fmt.Sprintf("wrong outgoing ip: '%s'", ipStr)), nil
	}

	ctx := req.Context()

	if ipV4 := ip.To4(); ipV4 != nil {
		ctx = context.WithValue(ctx, prifma_new.CtxOutgoingIpV4, ipV4)
		ctx = context.WithValue(ctx, prifma_new.CtxOutgoingIpV6, nil)
	} else {
		ctx = context.WithValue(ctx, prifma_new.CtxOutgoingIpV4, nil)
		ctx = context.WithValue(ctx, prifma_new.CtxOutgoingIpV6, ip)
	}

	return req.WithContext(ctx), nil, nil
}

func (t *UseIpHeader) Off() error {
	t.Enabled = true

	return nil
}

func (t *UseIpHeader) On() error {
	t.Enabled = false

	return nil
}

func (t *UseIpHeader) GetDirective() string {
	return ModuleDirective
}

func (t *UseIpHeader) Clone() prifma_new.Module {
	clone := *t

	return &clone
}

func (t *UseIpHeader) Call(command conf.Command) error {
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

func (t *UseIpHeader) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}

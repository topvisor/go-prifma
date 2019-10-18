package useipheader

import (
	"fmt"
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma"
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

func New() prifma.Module {
	return new(UseIpHeader)
}

func (t *UseIpHeader) HandleRequest(result prifma.HandleRequestResult) (prifma.HandleRequestResult, error) {
	if !t.Enabled {
		return result, nil
	}

	ipStr := result.GetRequest().Header.Get(HeaderName)
	if ipStr == "" {
		return result, nil
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		result.SetResponse(prifma.NewResponseError(http.StatusBadRequest, fmt.Sprintf("wrong outgoing ip: '%s'", ipStr)))

		return result, nil
	}

	result.GetDialer().SetIpV4(nil)
	result.GetDialer().SetIpV6(nil)

	if ipV4 := ip.To4(); ipV4 != nil {
		result.GetDialer().SetIpV4(ipV4)
	} else {
		result.GetDialer().SetIpV6(ip)
	}

	return result, nil
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

func (t *UseIpHeader) Clone() prifma.Module {
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

package proxyreq

import (
	"fmt"
	"github.com/topvisor/go-prifma/pkg/conf"
	"github.com/topvisor/go-prifma/pkg/prifma"
	"net/http"
	"net/url"
)

const ModuleDirective = "proxy_requests"

type UseIpHeader struct {
	Proxy       prifma.ProxyFunc
	ProxyHeader http.Header
}

func New() prifma.Module {
	return new(UseIpHeader)
}

func (t *UseIpHeader) HandleRequest(result prifma.HandleRequestResult) (prifma.HandleRequestResult, error) {
	if t.Proxy == nil {
		return result, nil
	}

	result.SetProxy(t.Proxy)
	result.SetProxyConnectHeader(t.ProxyHeader)

	return result, nil
}

func (t *UseIpHeader) Off() error {
	t.Proxy = nil
	t.ProxyHeader = nil

	return nil
}

func (t *UseIpHeader) SetProxyUrl(urlStr string) error {
	uri, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("wrong proxy url: '%s'", urlStr)
	}

	t.Proxy = http.ProxyURL(uri)
	t.ProxyHeader = make(map[string][]string)

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

	arg := command.GetArgs()[0]
	if arg == "off" {
		return t.Off()
	}

	return t.SetProxyUrl(arg)
}

func (t *UseIpHeader) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() != ModuleDirective || len(command.GetArgs()) != 1 {
		return nil, conf.NewErrCommand(command)
	}

	if err := t.SetProxyUrl(command.GetArgs()[0]); err != nil {
		return nil, err
	}

	return NewConfBlock(&t.ProxyHeader), nil
}

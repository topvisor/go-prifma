package proxyreq

import (
	"fmt"
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma_new"
	"net/http"
	"net/url"
)

const ModuleDirective = "proxy_requests"

type UseIpHeader struct {
	Proxy       prifma_new.ProxyFunc
	ProxyHeader http.Header
}

func New() prifma_new.Module {
	return new(UseIpHeader)
}

func (t *UseIpHeader) HandleRequest(result prifma_new.HandleRequestResult) (prifma_new.HandleRequestResult, error) {
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

func (t *UseIpHeader) Clone() prifma_new.Module {
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

	return nil, conf.NewErrCommand(command)
}
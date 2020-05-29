package http

import (
	"github.com/topvisor/go-prifma/pkg/conf"
	"github.com/topvisor/go-prifma/pkg/prifma"
	"github.com/topvisor/go-prifma/pkg/utils"
	"net/http"
)

const ModuleDirective = "http"

type Http struct {
	RoundTrippers RoundTrippersMap
}

func New() prifma.Module {
	return &Http{
		RoundTrippers: NewSyncRoundTrippersMap(),
	}
}

func (t *Http) HandleRequest(result prifma.HandleRequestResult) (prifma.HandleRequestResult, error) {
	if result.GetRequest().Method == http.MethodConnect {
		return result, nil
	}

	dialer := result.GetDialer()
	host := utils.GetRequestHostname(result.GetRequest())
	localIp, err := dialer.GetLocalIp(host)
	if err != nil {
		if err == prifma.ErrOutgoingIpNotDefined {
			result.SetResponse(prifma.NewResponseError(http.StatusBadRequest, err.Error()))
		} else {
			result.SetResponse(prifma.NewResponseError(http.StatusBadGateway, err.Error()))
		}

		return result, nil
	}

	if localIpV4 := localIp.To4(); localIpV4 != nil {
		dialer.SetIpV4(localIpV4)
		dialer.SetIpV6(nil)
	} else {
		dialer.SetIpV4(nil)
		dialer.SetIpV6(localIp)
	}

	result.SetResponse(NewResponseReverseProxy(t.RoundTrippers))

	return result, nil
}

func (t *Http) GetDirective() string {
	return ModuleDirective
}

func (t *Http) Clone() prifma.Module {
	return t
}

func (t *Http) Call(command conf.Command) error {
	return conf.NewErrCommandName(command)
}

func (t *Http) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommandName(command)
}

package http

import (
	"bytes"
	"github.com/topvisor/go-prifma/pkg/prifma"
	"github.com/topvisor/go-prifma/pkg/utils"
	"net/http"
	"sync"
)

type RoundTrippersMap interface {
	Get(result prifma.HandleRequestResult) http.RoundTripper
}

type SyncRoundTrippersMap struct {
	RWMutex       *sync.RWMutex
	RoundTrippers map[RoundTripperKey]http.RoundTripper
}

func NewSyncRoundTrippersMap() *SyncRoundTrippersMap {
	return &SyncRoundTrippersMap{
		RWMutex:       new(sync.RWMutex),
		RoundTrippers: make(map[RoundTripperKey]http.RoundTripper),
	}
}

func (t *SyncRoundTrippersMap) Get(result prifma.HandleRequestResult) http.RoundTripper {
	key := NewRoundTripperKey(result)

	t.RWMutex.RLock()
	roundTripper := t.RoundTrippers[key]
	t.RWMutex.RUnlock()

	if roundTripper != nil {
		return roundTripper
	}

	t.RWMutex.Lock()
	if roundTripper = t.RoundTrippers[key]; roundTripper == nil {
		roundTripper = result.GetRoundTripper()
		t.RoundTrippers[key] = roundTripper
	}
	t.RWMutex.Unlock()

	return roundTripper
}

type RoundTripperKey struct {
	ProxyUrl    string
	ProxyHeader string
	LocalIp     string
}

func NewRoundTripperKey(result prifma.HandleRequestResult) RoundTripperKey {
	t := RoundTripperKey{}

	if result.GetProxy() != nil {
		if proxyUrl, err := result.GetProxy()(result.GetRequest()); err == nil && proxyUrl != nil {
			t.ProxyUrl = proxyUrl.String()
		}
	}

	proxyHeaderBuff := new(bytes.Buffer)
	if err := result.GetProxyConnectHeader().Write(proxyHeaderBuff); err == nil {
		t.ProxyHeader = proxyHeaderBuff.String()
	}

	host := utils.GetHostname(result.GetRequest().Host)
	if localIp, err := result.GetDialer().GetLocalIp(host); err == nil && localIp != nil {
		t.LocalIp = localIp.String()
	}

	return t
}

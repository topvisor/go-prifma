package go_proxy_server

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Proxy http.Transport

func (t *Proxy) UnmarshalJSON(data []byte) error {
	var jsonProxy struct {
		Url            URL                `json:"url"`
		ConnectHeaders *map[string]string `json:"connectHeaders,omitempty"`
	}

	err := json.Unmarshal(data, &jsonProxy)
	if err != nil {
		return err
	}

	t.Proxy = http.ProxyURL((*url.URL)(&jsonProxy.Url))
	t.DialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	t.MaxIdleConns = 100
	t.IdleConnTimeout = 90 * time.Second
	t.TLSHandshakeTimeout = 10 * time.Second
	t.ExpectContinueTimeout = 1 * time.Second

	if jsonProxy.ConnectHeaders != nil {
		t.ProxyConnectHeader = http.Header{}

		for key, value := range *jsonProxy.ConnectHeaders {
			t.ProxyConnectHeader.Add(key, value)
		}
	}

	return nil
}

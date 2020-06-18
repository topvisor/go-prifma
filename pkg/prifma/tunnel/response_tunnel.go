package tunnel

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/topvisor/go-prifma/pkg/prifma"
	"github.com/topvisor/go-prifma/pkg/utils"
	"net"
	"net/http"
	"net/url"
)

type ResponseTunnel struct {
	ResponseCode int
	DstConn      net.Conn
}

func NewResponseTunnel() *ResponseTunnel {
	return new(ResponseTunnel)
}

func (t *ResponseTunnel) Write(rw http.ResponseWriter, result prifma.HandleRequestResult) error {
	if result.GetProxy() != nil {
		if err := t.ConnectToProxy(result); err != nil {
			http.Error(rw, err.Error(), http.StatusBadGateway)
			t.ResponseCode = http.StatusBadGateway

			return err
		}
	}

	if t.DstConn == nil {
		if err := t.ConnectToRequest(result); err != nil {
			http.Error(rw, err.Error(), http.StatusBadGateway)
			t.ResponseCode = http.StatusBadGateway

			return err
		}
	}

	rw.WriteHeader(http.StatusOK)
	t.ResponseCode = http.StatusOK

	clientConn, _, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		utils.CloseFile(t.DstConn)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		t.ResponseCode = http.StatusInternalServerError

		return err
	}

	writeTimeout := result.GetServer().GetWriteTimeout()
	readTimeout := result.GetServer().GetReadTimeout()
	if readTimeout == 0 {
		readTimeout = result.GetServer().GetReadHeaderTimeout()
	}
	readTimeout += result.GetServer().GetIdleTimeout()

	go utils.Transfer(readTimeout, writeTimeout, clientConn, t.DstConn)
	go utils.Transfer(readTimeout, writeTimeout, t.DstConn, clientConn)

	return nil
}

func (t *ResponseTunnel) GetCode() int {
	return t.ResponseCode
}

func (t *ResponseTunnel) GetLAddr() net.Addr {
	if t.DstConn == nil {
		return nil
	}

	return t.DstConn.LocalAddr()
}

func (t *ResponseTunnel) GetRAddr() net.Addr {
	if t.DstConn == nil {
		return nil
	}

	return t.DstConn.RemoteAddr()
}

func (t *ResponseTunnel) ConnectToProxy(result prifma.HandleRequestResult) error {
	proxyUrl, err := result.GetProxy()(result.GetRequest())
	if err != nil {
		return err
	}

	req := &http.Request{
		Method: http.MethodConnect,
		URL: &url.URL{
			Scheme: proxyUrl.Scheme,
			Host:   result.GetRequest().Host,
		},
		Header: make(http.Header),
		Host:   result.GetRequest().Host,
	}

	req.Header.Set("Host", result.GetRequest().Host)
	req.Header.Set("Proxy-Connection", "keep-alive")
	if proxyUrl.User != nil {
		authHash := base64.StdEncoding.EncodeToString([]byte(proxyUrl.User.String()))
		req.Header.Set("Proxy-Authorization", "Basic "+authHash)
	}
	for key, values := range result.GetProxyConnectHeader() {
		req.Header[key] = values
	}

	req = req.WithContext(result.GetRequest().Context())

	roundTripper := &http.Transport{
		DialContext: func(ctx context.Context, network, _ string) (conn net.Conn, err error) {
			t.DstConn, err = result.GetDialer().DialContext(ctx, network, proxyUrl.Host)

			return t.DstConn, err
		},
		ResponseHeaderTimeout: result.GetServer().GetWriteTimeout(),
		TLSNextProto:          make(map[string]func(string, *tls.Conn) http.RoundTripper, 0), // disable HTTP/2
	}

	resp, err := roundTripper.RoundTrip(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(http.StatusBadGateway))
	}

	return nil
}

func (t *ResponseTunnel) ConnectToRequest(result prifma.HandleRequestResult) error {
	host, port, err := net.SplitHostPort(result.GetRequest().Host)
	if err != nil {
		host = result.GetRequest().Host
		port = "443"
	}

	ctx := result.GetRequest().Context()
	network := "tcp"
	addr := net.JoinHostPort(host, port)

	t.DstConn, err = result.GetDialer().DialContext(ctx, network, addr)

	return err
}

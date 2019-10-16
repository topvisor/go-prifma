package prifma

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

type Proxy struct {
	Url          *url.URL
	ProxyHeaders http.Header
}

func NewProxyFromConfig(config ConfigProxy) (*Proxy, error) {
	p := new(Proxy)
	if err := p.SetFromConfig(config); err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Proxy) SetFromConfig(config ConfigProxy) error {
	parserUrl, err := url.ParseRequestURI(config.Url)
	if err != nil {
		return err
	}

	headers := http.Header{}
	if config.ProxyHeaders != nil {
		for key, val := range config.ProxyHeaders {
			headers.Add(key, val)
		}
	}

	t.Url = parserUrl
	t.ProxyHeaders = headers

	return nil
}

func (t *Proxy) connect(req *http.Request, dialer *dialer) (net.Conn, error) {
	conn, err := dialer.connect(t.Url)
	if err != nil {
		return nil, err
	}

	connectRequest := fmt.Sprintf("CONNECT %s HTTP/1.1\r\n", req.Host)
	connectRequest += fmt.Sprintf("Host: %s\r\n", req.Host)
	connectRequest += "Proxy-Connection: keep-alive\r\n"
	if t.Url.User != nil {
		authHash := base64.StdEncoding.EncodeToString([]byte(t.Url.User.String()))
		connectRequest += fmt.Sprintf("Proxy-Authorization: Basic %s\r\n", authHash)
	}
	if len(t.ProxyHeaders) != 0 {
		buf := new(bytes.Buffer)
		err := t.ProxyHeaders.Write(buf)
		if err != nil {
			return nil, err
		}

		connectRequest += buf.String()
	}
	connectRequest += "\r\n"

	_, err = conn.Write([]byte(connectRequest))
	if err != nil {
		return nil, err
	}

	respReader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(respReader, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		bodyData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		body := string(bodyData)
		err = fmt.Errorf("CONNECT StatusCode: %d\n%s", resp.StatusCode, body)

		return nil, err
	}

	return conn, nil
}

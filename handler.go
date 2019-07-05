package proxy

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Handler httputil.ReverseProxy

func NewHandler(filters *[]Filter) *Handler {
	var handler Handler
	handler.Director = handler.director

	if filters != nil {
		handler.Transport = (*http.Transport)((*filters)[0].Proxy)
	}

	return &handler
}

func (t *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		t.serveTunnel(rw, req)
	} else {
		(*httputil.ReverseProxy)(t).ServeHTTP(rw, req)
	}
}

func (t *Handler) serveTunnel(rw http.ResponseWriter, req *http.Request) {
	var destConn net.Conn

	if t.Transport != nil {
		if proxyUrl, _ := t.Transport.(*http.Transport).Proxy(nil); proxyUrl != nil {
			var err error
			if destConn, err = t.connectToProxy(req, proxyUrl); err != nil {
				http.Error(rw, err.Error(), http.StatusServiceUnavailable)
				return
			}
		}
	}

	if destConn == nil {
		var err error
		if destConn, err = connectToHost(req.Host); err != nil {
			http.Error(rw, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)

	clientConn, _, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		_ = destConn.Close()
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}

	go transfer(clientConn, destConn)
	go transfer(destConn, clientConn)
}

func (t *Handler) director(req *http.Request) {
	log.Print(req)
}

func (t *Handler) connectToProxy(req *http.Request, proxyUrl *url.URL) (net.Conn, error) {
	proxyHost := proxyUrl.Host
	if proxyUrl.Port() == "" {
		proxyHost = net.JoinHostPort(proxyHost, "80")
	}

	conn, err := connectToHost(proxyHost)
	if err != nil {
		return nil, err
	}

	proxyHeader := t.Transport.(*http.Transport).ProxyConnectHeader

	connectRequest := fmt.Sprintf("CONNECT %s HTTP/1.1\r\n", req.Host)
	connectRequest += fmt.Sprintf("Host: %s\r\n", req.Host)
	connectRequest += "Proxy-Connection: keep-alive\r\n"
	if proxyUrl.User != nil {
		authHash := base64.StdEncoding.EncodeToString([]byte(proxyUrl.User.String()))
		connectRequest += fmt.Sprintf("Proxy-Authorization: Basic %s\r\n", authHash)
	}
	if len(proxyHeader) != 0 {
		buf := new(bytes.Buffer)
		err := proxyHeader.Write(buf)
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

func connectToHost(host string) (net.Conn, error) {
	return net.DialTimeout("tcp", host, 10*time.Second)
}

func transfer(src io.ReadCloser, dst io.WriteCloser) {
	_, _ = io.Copy(dst, src)
	_ = src.Close()
	_ = dst.Close()
}

package proxy

import (
	"context"
	"errors"
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// condWithHandler contains the condition by which the handler will be selected
type condWithHandler struct {
	tester  condition
	handler *Handler
}

// conditionUniqueKey is a key which identify a condition
type conditionUniqueKey struct {
	Type  conditionType
	Value string
}

const keyReqId = iota + 413

const (
	StatusClientClosedRequest     = 499
	StatusTextClientClosedRequest = "Client Closed Request"
)

// Handler
type Handler struct {
	AccessLog         *log.Logger
	HandleTimeout     *time.Duration
	BasicAuth         *auth.BasicAuth
	EnableBasicAuth   *bool
	OutgoingIpV4      []net.IP
	OutgoingIpV6      []net.IP
	EnableUseIpHeader *bool
	BlockRequests     *bool
	Proxy             *Proxy

	server         *server
	conditions     map[conditionUniqueKey]*condWithHandler
	reverseProxies sync.Map
	nextReqId      uint64
}

func (t *Handler) SetConditionHandler(cond *Condition, handler Handler) {
	if t.conditions == nil {
		t.conditions = make(map[conditionUniqueKey]*condWithHandler)
	}

	condUniqueKey := conditionUniqueKey{cond.Type, cond.Value}

	if handler.AccessLog == nil {
		handler.AccessLog = t.AccessLog
	}
	if handler.HandleTimeout == nil {
		handler.HandleTimeout = t.HandleTimeout
	}
	if handler.BasicAuth == nil {
		handler.BasicAuth = t.BasicAuth
	}
	if handler.EnableBasicAuth == nil {
		handler.EnableBasicAuth = t.EnableBasicAuth
	}
	if handler.OutgoingIpV4 == nil {
		handler.OutgoingIpV4 = t.OutgoingIpV4
	}
	if handler.OutgoingIpV6 == nil {
		handler.OutgoingIpV6 = t.OutgoingIpV6
	}
	if handler.EnableUseIpHeader == nil {
		handler.EnableUseIpHeader = t.EnableUseIpHeader
	}
	if handler.BlockRequests == nil {
		handler.BlockRequests = t.BlockRequests
	}
	if handler.Proxy == nil {
		handler.Proxy = t.Proxy
	}

	t.conditions[condUniqueKey] = &condWithHandler{cond, &handler}
}

func (t *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handler := t.getHandler(req)
	handler.serveHTTP(rw, req)
}

func (t *Handler) setServer(server *server) {
	t.server = server

	if t.conditions != nil {
		for _, condWithHandler := range t.conditions {
			condWithHandler.handler.setServer(server)
		}
	}
}

func (t *Handler) setFromConfig(config ConfigHandler) error {
	if config.AccessLog != nil {
		accessLogFile, err := os.Create(*config.AccessLog)
		if err != nil {
			return err
		}

		t.AccessLog = log.New(accessLogFile, "", log.Flags())
	}
	if config.HandleTimeout != nil {
		handleTimeout, err := time.ParseDuration(*config.HandleTimeout)
		if err != nil {
			return err
		}

		t.HandleTimeout = &handleTimeout
	}
	if config.Htpasswd != nil {
		htpasswd, err := LoadHtpasswd(*config.Htpasswd)
		if err != nil {
			return err
		}

		t.BasicAuth = &htpasswd.BasicAuth
	}
	if t.EnableBasicAuth = config.EnableBasicAuth; t.EnableBasicAuth != nil && *t.EnableBasicAuth && t.BasicAuth == nil {
		return errors.New(".BasicAuth must be set")
	}
	if config.OutgoingIpV4 != nil {
		t.OutgoingIpV4 = make([]net.IP, len(config.OutgoingIpV4.Ips))
		for i, ip := range config.OutgoingIpV4.Ips {
			if t.OutgoingIpV4[i] = net.ParseIP(ip); t.OutgoingIpV4[i] == nil || strings.Contains(ip, ":") {
				return fmt.Errorf("incorrect outgoing ip v4 address(es): %v", config.OutgoingIpV4)
			}
		}
	}
	if config.OutgoingIpV6 != nil {
		t.OutgoingIpV6 = make([]net.IP, len(config.OutgoingIpV6.Ips))
		for i, ip := range config.OutgoingIpV6.Ips {
			if t.OutgoingIpV6[i] = net.ParseIP(ip); t.OutgoingIpV6[i] == nil || strings.Contains(ip, ".") {
				return fmt.Errorf("incorrect outgoing ip v6 address(es): %v", config.OutgoingIpV6)
			}
		}
	}

	t.EnableUseIpHeader = config.EnableUseIpHeader
	t.BlockRequests = config.BlockRequests

	if config.Proxy != nil {
		t.Proxy = new(Proxy)
		if err := t.Proxy.SetFromConfig(*config.Proxy); err != nil {
			return err
		}
	}
	if config.Conditions != nil {
		for _, configCondition := range config.Conditions {
			condition, err := parseConditionFromString(configCondition.Condition)
			if err != nil {
				return err
			}

			handler := Handler{}
			if err = handler.setFromConfig(configCondition.Handler); err != nil {
				return err
			}

			t.SetConditionHandler(condition, handler)
		}
	}

	return nil
}

func (t *Handler) getHandler(req *http.Request) *Handler {
	if t.conditions == nil {
		return t
	}

	for _, condAndHandler := range t.conditions {
		if condAndHandler.tester.test(req) {
			return condAndHandler.handler
		}
	}

	return t
}

func (t *Handler) serveHTTP(rw http.ResponseWriter, req *http.Request) {
	reqId := atomic.AddUint64(&t.nextReqId, 1)

	ctx := req.Context()
	ctx = context.WithValue(ctx, keyReqId, reqId)
	if t.HandleTimeout != nil {
		ctx, _ = context.WithTimeout(ctx, *t.HandleTimeout)
	}

	req = req.WithContext(ctx)
	respWriterChan := make(chan responseWriter)

	go func() {
		respWriterChan <- t.serveHTTPContext(req)
		close(respWriterChan)
	}()

	var respWriter responseWriter
	select {
	case respWriter = <-respWriterChan:
	case <-ctx.Done():
	}

	switch ctx.Err() {
	case context.DeadlineExceeded:
		respWriter = &responseWriterError{Code: http.StatusGatewayTimeout}
	case context.Canceled:
		respWriter = &responseWriterError{Code: StatusClientClosedRequest, Error: StatusTextClientClosedRequest}
	}

	if respWriter == nil {
		respWriter = &responseWriterError{Code: http.StatusInternalServerError}
	}
	if err := respWriter.Write(rw); err != nil {
		t.server.ErrorLog.Println(err)
	}

	if t.AccessLog != nil {
		var user *string
		if username, _, ok := req.BasicAuth(); ok {
			user = &username
		}

		t.AccessLog.Printf(
			"%s %d %s %s %v l/%v r/%v\n",
			req.RemoteAddr,
			respWriter.GetCode(),
			req.Method,
			req.RequestURI,
			user,
			respWriter.GetLAddr(),
			respWriter.GetRAddr(),
		)
	}
	// ### access log
}

func (t *Handler) serveHTTPContext(req *http.Request) responseWriter {
	if t.EnableBasicAuth != nil && *t.EnableBasicAuth && t.BasicAuth != nil && t.BasicAuth.CheckAuth(req) == "" {
		return &responseWriteRequireAuth{req, t.BasicAuth}
	}
	if t.BlockRequests != nil && *t.BlockRequests {
		return &responseWriterError{Code: http.StatusLocked}
	}

	if req.Method == http.MethodConnect {
		return t.serveTunnel(req)
	} else {
		return t.serveReverseProxy(req)
	}
}

func (t *Handler) serveTunnel(req *http.Request) responseWriter {
	var destConn net.Conn
	var err error
	var dialer *dialer

	if dialer, err = t.getDialer(req); err != nil {
		return &responseWriterError{Code: http.StatusBadGateway}
	}
	if t.Proxy != nil {
		if destConn, err = t.Proxy.connect(req, dialer); err != nil {
			return &responseWriterError{Code: http.StatusBadGateway}
		}
	}
	if destConn == nil {
		destUrl := &url.URL{Host: req.Host}
		if destUrl.Port() == "" {
			destUrl.Host = net.JoinHostPort(req.Host, "443")
		}
		if destConn, err = dialer.connect(destUrl); err != nil {
			return &responseWriterError{Code: http.StatusBadGateway}
		}
	}

	return &responseWriterTunnel{
		DestConn:     destConn,
		ReadTimeout:  t.server.ReadTimeout,
		WriteTimeout: t.server.WriteTimeout,
	}
}

func (t *Handler) serveReverseProxy(req *http.Request) responseWriter {
	dialer, err := t.getDialer(req)
	if err != nil {
		return &responseWriterError{Code: http.StatusBadGateway}
	}

	ipsKey := dialer.ipsString()

	var rProxy *reverseProxy
	if val, ok := t.reverseProxies.Load(ipsKey); ok {
		rProxy = val.(*reverseProxy)
	} else {
		rProxy = newReverseProxy(t, dialer)
		t.reverseProxies.Store(ipsKey, rProxy)
	}

	return &responseWriteReverseProxy{
		Request:      req,
		ReverseProxy: rProxy,
	}
}

func (t *Handler) getDialer(req *http.Request) (*dialer, error) {
	dialer := new(dialer)

	if t.OutgoingIpV4 != nil {
		dialer.lIpV4 = t.OutgoingIpV4[rand.Intn(len(t.OutgoingIpV4))]
	}
	if t.OutgoingIpV6 != nil {
		dialer.lIpV6 = t.OutgoingIpV6[rand.Intn(len(t.OutgoingIpV6))]
	}

	if t.EnableUseIpHeader != nil && *t.EnableUseIpHeader {
		if lIpStr := req.Header.Get("Proxy-Use-IpV4"); lIpStr != "" {
			lIp := net.ParseIP(lIpStr)
			if lIp == nil || strings.Contains(lIpStr, ":") {
				return nil, fmt.Errorf("incorrect outgoing ip v4 address: \"%s\"", lIpStr)
			}

			dialer.lIpV4 = lIp
		}
		if lIpStr := req.Header.Get("Proxy-Use-IpV6"); lIpStr != "" {
			lIp := net.ParseIP(lIpStr)
			if lIp == nil || strings.Contains(lIpStr, ".") {
				return nil, fmt.Errorf("incorrect outgoing ip v6 address: \"%s\"", lIpStr)
			}

			dialer.lIpV6 = lIp
		}
	}

	return dialer, nil
}

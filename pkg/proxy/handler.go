package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type condWithHandler struct {
	tester  condition
	handler *Handler
}

type conditionUniqueKey struct {
	Type  conditionType
	Value string
}

type Handler struct {
	AccessLogger      Logger
	ErrorLogger       Logger
	DialTimeout       *time.Duration
	Htpasswd          *BasicAuth
	EnableBasicAuth   *bool
	OutgoingIpV4      net.IP
	OutgoingIpV6      net.IP
	EnableUseIpHeader *bool
	BlockRequests     *bool
	Proxy             *Proxy

	conditions map[conditionUniqueKey]*condWithHandler
}

func (t *Handler) Close() error {
	if err := t.AccessLogger.Close(); err != nil {
		return err
	}
	if err := t.ErrorLogger.Close(); err != nil {
		return err
	}

	for _, condWithHandler := range t.conditions {
		if err := condWithHandler.handler.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (t *Handler) SetConditionHandler(cond *Condition, handler Handler) error {
	if t.conditions == nil {
		t.conditions = make(map[conditionUniqueKey]*condWithHandler)
	}

	condUniqueKey := conditionUniqueKey{cond.Type, cond.Value}

	if oldCondAndHandler, exists := t.conditions[condUniqueKey]; exists {
		if err := oldCondAndHandler.handler.Close(); err != nil {
			return err
		}
	}

	if !handler.AccessLogger.IsInited() {
		handler.AccessLogger = t.AccessLogger
	}
	if !handler.ErrorLogger.IsInited() {
		handler.ErrorLogger = t.ErrorLogger
	}
	if handler.DialTimeout == nil {
		handler.DialTimeout = t.DialTimeout
	}
	if handler.Htpasswd == nil {
		handler.Htpasswd = t.Htpasswd
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

	return nil
}

func (t *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handler := t.getHandler(req)
	handler.serveHTTP(rw, req)
}

func (t *Handler) setFromConfig(config ConfigHandler) error {
	var err error

	if config.AccessLog != nil {
		if err = t.AccessLogger.SetFile(*config.AccessLog); err != nil {
			return err
		}
	}
	if config.ErrorLog != nil {
		if err = t.ErrorLogger.SetFile(*config.ErrorLog); err != nil {
			return err
		}
	}
	if config.DialTimeout == nil {
		dialTimeout := time.Second * time.Duration(*config.DialTimeout)
		t.DialTimeout = &dialTimeout
	}
	if config.Htpasswd != nil {
		if t.Htpasswd, err = NewBasicAuth(*config.Htpasswd); err != nil {
			return err
		}
	}
	if t.EnableBasicAuth = config.EnableBasicAuth; t.EnableBasicAuth != nil && *t.EnableBasicAuth && t.Htpasswd == nil {
		return errors.New(".htpasswd must be set")
	}
	if config.OutgoingIpV4 != nil {
		if t.OutgoingIpV4 = net.ParseIP(*config.OutgoingIpV4); t.OutgoingIpV4 == nil || strings.Contains(*config.OutgoingIpV4, ":") {
			return fmt.Errorf("incorrect outgoing ip v4 address: \"%s\"", *config.OutgoingIpV4)
		}
	}
	if config.OutgoingIpV6 != nil {
		if t.OutgoingIpV6 = net.ParseIP(*config.OutgoingIpV6); t.OutgoingIpV6 == nil || strings.Contains(*config.OutgoingIpV6, ".") {
			return fmt.Errorf("incorrect outgoing ip v6 address: \"%s\"", *config.OutgoingIpV6)
		}
	}

	t.EnableUseIpHeader = config.EnableUseIpHeader
	t.BlockRequests = config.BlockRequests

	if config.Proxy != nil {
		t.Proxy = new(Proxy)
		if err = t.Proxy.setFromConfig(*config.Proxy); err != nil {
			return err
		}
	}
	if config.Conditions != nil {
		for _, configCondition := range config.Conditions {
			condition, err := ParseConditionFromString(configCondition.Condition)
			if err != nil {
				return err
			}

			handler := Handler{}
			if err = handler.setFromConfig(configCondition.Handler); err != nil {
				return err
			}

			if err = t.SetConditionHandler(condition, handler); err != nil {
				return nil
			}
		}
	}

	return nil
}

func (t *Handler) getHandler(req *http.Request) *Handler {
	if t.conditions == nil {
		return t
	}

	for _, condAndHandler := range t.conditions {
		if condAndHandler.tester.Test(req) {
			return condAndHandler.handler
		}
	}

	return t
}

func (t *Handler) serveHTTP(rw http.ResponseWriter, req *http.Request) {
	if t.EnableBasicAuth != nil && *t.EnableBasicAuth && t.Htpasswd != nil && t.Htpasswd.CheckAuth(req) == "" {
		t.Htpasswd.RequireAuth(rw, req)
		return
	}
	if t.BlockRequests != nil && *t.BlockRequests {
		http.Error(rw, http.StatusText(http.StatusLocked), http.StatusLocked)
		return
	}

	if req.Method == http.MethodConnect {
		t.serveTunnel(rw, req)
	} else {
		//(*httputil.ReverseProxy)(t).ServeHTTP(rw, req)
	}
}

func (t *Handler) serveTunnel(rw http.ResponseWriter, req *http.Request) {
	var destConn net.Conn
	var err error
	var dialer *dialer

	if dialer, err = t.getDialer(req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		//t.AccessLogger.Println() // ###
		return
	}
	if t.Proxy != nil {
		if destConn, err = t.Proxy.connect(req, dialer); err != nil {
			http.Error(rw, err.Error(), http.StatusBadGateway)
			//t.AccessLogger.Println() // ###
			return
		}
	}
	if destConn == nil {
		destUrl := &url.URL{Host: req.Host}
		if destUrl.Port() == "" {
			destUrl.Host = net.JoinHostPort(req.Host, "443")
		}
		if destConn, err = dialer.connect(destUrl); err != nil {
			http.Error(rw, err.Error(), http.StatusBadGateway)
			//t.AccessLogger.Println() // ###
			return
		}
	}

	rw.WriteHeader(http.StatusOK)

	clientConn, _, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		_ = destConn.Close()
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		//t.ErrorLogger.Println() // ###
		return
	}

	//t.AccessLogger.Println() // ###

	go transfer(clientConn, destConn)
	go transfer(destConn, clientConn)
}

func (t *Handler) getDialer(req *http.Request) (*dialer, error) {
	dialer := new(dialer)

	dialer.lIpV4 = t.OutgoingIpV4
	dialer.lIpV6 = t.OutgoingIpV6

	if t.DialTimeout != nil {
		dialer.Timeout = *t.DialTimeout
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

func transfer(src io.ReadCloser, dst io.WriteCloser) {
	_, _ = io.Copy(dst, src)
	_ = src.Close()
	_ = dst.Close()
}

package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type listenType byte

const (
	ListenTypeHttp listenType = iota
)

func ListenTypeFromString(lTypeStr string) (*listenType, error) {
	switch lTypeStr {
	case "http":
		listenType := ListenTypeHttp
		return &listenType, nil
	default:
		return nil, fmt.Errorf("unavailable listen type: \"%s\"", lTypeStr)
	}
}

type Server struct {
	ListenIp   net.IP
	ListenPort int
	ListenType listenType
	Handler    Handler

	httpServer http.Server
}

func (t *Server) SetFromConfig(config Config) error {
	port := config.Listen.ListenPort

	var ip net.IP
	if config.Listen.ListenIp != nil {
		if ip = net.ParseIP(*config.Listen.ListenIp); ip == nil {
			return fmt.Errorf("incorrect ip address: \"%s\"", *config.Listen.ListenIp)
		}
	}

	ltype, err := ListenTypeFromString(config.Listen.ListenType)
	if err != nil {
		return err
	}

	handler := Handler{}
	if err = handler.setFromConfig(config.ConfigHandler); err != nil {
		return err
	}
	if err = t.Handler.Close(); err != nil {
		return err
	}

	t.ListenIp = ip
	t.ListenPort = port
	t.ListenType = *ltype
	t.Handler = handler

	return nil
}

func (t *Server) LoadFromConfig(filename string) error {
	config, err := ParseConfigFromFile(filename)
	if err != nil {
		return err
	}

	if err = t.SetFromConfig(*config); err != nil {
		return err
	}

	return nil
}

func (t *Server) ListenAndServe() error {
	defer func() {
		_ = t.Handler.Close()
	}()

	var err error

	ipStr := ""
	if t.ListenIp != nil {
		ipStr = t.ListenIp.String()
	}

	t.httpServer.Addr = fmt.Sprintf("%s:%d", ipStr, t.ListenPort)
	t.httpServer.Handler = &t.Handler

	switch t.ListenType {
	case ListenTypeHttp:
		if err = t.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			t.Handler.ErrorLogger.Println(err)
		}
	default:
		err = fmt.Errorf("unavailable listen type: \"%v\"", t.ListenType)
	}

	return err
}

func (t *Server) Close() error {
	defer func() {
		_ = t.Handler.Close()
	}()

	return t.httpServer.Close()
}

func (t *Server) Shutdown(ctx context.Context) error {
	defer func() {
		_ = t.Handler.Close()
	}()

	return t.httpServer.Shutdown(ctx)
}

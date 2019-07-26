package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type listenType byte

const (
	ListenTypeHttp listenType = iota
)

func listenTypeFromString(lTypeStr string) (*listenType, error) {
	switch lTypeStr {
	case "http":
		listenType := ListenTypeHttp
		return &listenType, nil
	default:
		return nil, fmt.Errorf("unavailable listen type: \"%s\"", lTypeStr)
	}
}

type Server struct {
	ListenIp          net.IP
	ListenPort        int
	ListenType        listenType
	ErrorLogger       Logger
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	Handler           Handler

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

	ltype, err := listenTypeFromString(config.Listen.ListenType)
	if err != nil {
		return err
	}

	var errorLogger Logger
	if config.Listen.ErrorLog != nil {
		if err = errorLogger.SetFile(*config.Listen.ErrorLog); err != nil {
			return err
		}
	}

	var readTimeout, readHeaderTimeout, writeTimeout, idleTimeout time.Duration
	if config.Listen.ReadTimeout != nil {
		if readTimeout, err = time.ParseDuration(*config.Listen.ReadTimeout); err != nil {
			return err
		}
	}
	if config.Listen.ReadHeaderTimeout != nil {
		if readHeaderTimeout, err = time.ParseDuration(*config.Listen.ReadHeaderTimeout); err != nil {
			return err
		}
	}
	if config.Listen.WriteTimeout != nil {
		if writeTimeout, err = time.ParseDuration(*config.Listen.WriteTimeout); err != nil {
			return err
		}
	}
	if config.Listen.IdleTimeout != nil {
		if idleTimeout, err = time.ParseDuration(*config.Listen.IdleTimeout); err != nil {
			return err
		}
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
	t.ErrorLogger = errorLogger
	t.ReadTimeout = readTimeout
	t.ReadHeaderTimeout = readHeaderTimeout
	t.WriteTimeout = writeTimeout
	t.IdleTimeout = idleTimeout
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
		_ = t.ErrorLogger.Close()
	}()

	var err error

	t.Handler.setErrorLogger(&t.ErrorLogger)

	ipStr := ""
	if t.ListenIp != nil {
		ipStr = t.ListenIp.String()
	}

	t.httpServer.Addr = fmt.Sprintf("%s:%d", ipStr, t.ListenPort)
	t.httpServer.Handler = &t.Handler
	t.httpServer.ErrorLog = t.ErrorLogger.logger
	t.httpServer.ReadTimeout = t.ReadTimeout
	t.httpServer.ReadHeaderTimeout = t.ReadHeaderTimeout
	t.httpServer.WriteTimeout = t.WriteTimeout
	t.httpServer.IdleTimeout = t.IdleTimeout

	switch t.ListenType {
	case ListenTypeHttp:
		if err = t.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			t.ErrorLogger.Println(err)
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

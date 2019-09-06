package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
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

type Server interface {
	ListenAndServe() error
	Close() error
	Shutdown(ctx context.Context) error
}

func NewServerFromConfig(config Config) (Server, error) {
	b := new(ServerBuilder)
	if err := b.SetFromConfig(config); err != nil {
		return nil, err
	}

	s := b.Build()

	return s, nil
}

func NewServerFromConfigFile(filename string) (Server, error) {
	b := new(ServerBuilder)
	if err := b.SetFromConfigFile(filename); err != nil {
		return nil, err
	}

	s := b.Build()

	return s, nil
}

type ServerBuilder struct {
	ListenIp          net.IP
	ListenPort        int
	ListenType        listenType
	ErrorLog          *log.Logger
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	Handler           Handler

	errorLoggerCloser io.Closer
}

func (t *ServerBuilder) SetFromConfig(config Config) error {
	port := config.Listen.ListenPort

	var ip net.IP
	if config.Listen.ListenIp != nil {
		if ip = net.ParseIP(*config.Listen.ListenIp); ip == nil {
			return fmt.Errorf("incorrect ip address: \"%s\"", *config.Listen.ListenIp)
		}
	}

	listenType, err := listenTypeFromString(config.Listen.ListenType)
	if err != nil {
		return err
	}

	var errorLog *log.Logger
	if config.Listen.ErrorLog != nil {
		errorLogFile, err := os.Create(*config.Listen.ErrorLog)
		if err != nil {
			return err
		}

		errorLog = log.New(errorLogFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
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

	t.ListenIp = ip
	t.ListenPort = port
	t.ListenType = *listenType
	t.ErrorLog = errorLog
	t.ReadTimeout = readTimeout
	t.ReadHeaderTimeout = readHeaderTimeout
	t.WriteTimeout = writeTimeout
	t.IdleTimeout = idleTimeout
	t.Handler = handler

	return nil
}

func (t *ServerBuilder) SetFromConfigFile(filename string) error {
	config, err := ParseConfigFromFile(filename)
	if err != nil {
		return err
	}

	return t.SetFromConfig(*config)
}

func (t *ServerBuilder) Build() Server {
	ipStr := ""
	if t.ListenIp != nil {
		ipStr = t.ListenIp.String()
	}

	addr := fmt.Sprintf("%s:%d", ipStr, t.ListenPort)

	errorLog := t.ErrorLog
	if errorLog == nil {
		errorLog = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	s := &server{
		ListenIp:          t.ListenIp,
		ListenPort:        t.ListenPort,
		ListenType:        t.ListenType,
		ErrorLog:          errorLog,
		ReadTimeout:       t.ReadTimeout,
		ReadHeaderTimeout: t.ReadHeaderTimeout,
		WriteTimeout:      t.WriteTimeout,
		IdleTimeout:       t.IdleTimeout,
		Handler:           t.Handler,

		Server: http.Server{
			Addr:              addr,
			Handler:           &t.Handler,
			ErrorLog:          errorLog,
			ReadTimeout:       t.ReadTimeout,
			ReadHeaderTimeout: t.ReadHeaderTimeout,
			WriteTimeout:      t.WriteTimeout,
			IdleTimeout:       t.IdleTimeout,
		},
	}

	t.Handler.setServer(s)

	return s
}

type server struct {
	ListenIp          net.IP
	ListenPort        int
	ListenType        listenType
	ErrorLog          *log.Logger
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	Handler           Handler

	http.Server
}

func (t *server) ListenAndServe() (err error) {
	switch t.ListenType {
	case ListenTypeHttp:
		if err = t.Server.ListenAndServe(); err != http.ErrServerClosed {
			t.ErrorLog.Println(err)
		}
	default:
		err = fmt.Errorf("unavailable listen type: \"%v\"", t.ListenType)
	}

	return err
}

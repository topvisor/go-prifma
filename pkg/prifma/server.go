package prifma

import (
	"crypto/tls"
	"fmt"
	"github.com/topvisor/go-prifma/pkg/conf"
	"github.com/topvisor/go-prifma/pkg/utils"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server interface {
	GetModulesManager() ModulesManager
	GetListenIp() net.IP
	GetListenPort() int
	GetListenType() ListenType
	GetCertFile() string
	GetKeyFile() string
	GetErrorLog() *log.Logger
	GetDebugLog() *log.Logger
	GetReadTimeout() time.Duration
	GetReadHeaderTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetIdleTimeout() time.Duration

	SetListenIp(ip string) error
	SetListenPort(port string) error
	SetListenType(typ string) error
	SetCertFile(filename string)
	SetKeyFile(filename string)
	SetErrorLog(filename string) error
	SetDebugLog(filename string) error
	SetReadTimeout(timeout string) error
	SetReadHeaderTimeout(timeout string) error
	SetWriteTimeout(timeout string) error
	SetIdleTimeout(timeout string) error

	LoadConfig(filename string) error
	ListenAndServe() error
}

func NewServer(modules ...Module) *DefaultServer {
	t := &DefaultServer{
		ModulesManager: NewModulesManager(modules...),
		ListenType:     ListenTypeHttp,
		ErrorLog:       log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds),
	}

	t.Config = NewConfigMain(t)
	t.Server.Handler = NewRequestHandler(t)
	t.Server.Addr = net.JoinHostPort("0.0.0.0", "3128")
	t.Server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0) // disable HTTP/2

	return t
}

type DefaultServer struct {
	ModulesManager ModulesManager
	ListenType     ListenType
	ErrorLog       *log.Logger
	DebugLog       *log.Logger
	CertFile       string
	KeyFile        string
	Config         conf.Block
	Server         http.Server
}

func (t *DefaultServer) GetModulesManager() ModulesManager {
	return t.ModulesManager
}

func (t *DefaultServer) GetListenIp() net.IP {
	ipStr, _, _ := net.SplitHostPort(t.Server.Addr)
	ip := net.ParseIP(ipStr)

	return ip
}

func (t *DefaultServer) GetListenPort() int {
	_, portStr, _ := net.SplitHostPort(t.Server.Addr)
	port, _ := strconv.Atoi(portStr)

	return port
}

func (t *DefaultServer) GetListenType() ListenType {
	return t.ListenType
}

func (t *DefaultServer) GetCertFile() string {
	return t.CertFile
}

func (t *DefaultServer) GetKeyFile() string {
	return t.KeyFile
}

func (t *DefaultServer) GetErrorLog() *log.Logger {
	return t.ErrorLog
}

func (t *DefaultServer) GetDebugLog() *log.Logger {
	return t.DebugLog
}

func (t *DefaultServer) GetReadTimeout() time.Duration {
	return t.Server.ReadTimeout
}

func (t *DefaultServer) GetReadHeaderTimeout() time.Duration {
	return t.Server.ReadHeaderTimeout
}

func (t *DefaultServer) GetWriteTimeout() time.Duration {
	return t.Server.WriteTimeout
}

func (t *DefaultServer) GetIdleTimeout() time.Duration {
	return t.Server.IdleTimeout
}

func (t *DefaultServer) SetListenIp(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid ip - %s", ip)
	}

	_, port, _ := net.SplitHostPort(t.Server.Addr)
	t.Server.Addr = net.JoinHostPort(ip, port)

	return nil
}

func (t *DefaultServer) SetListenPort(port string) error {
	if intPort, err := strconv.ParseUint(port, 0, 0); err != nil || intPort < 1 || intPort > 65535 {
		return fmt.Errorf("invalid port - %s", port)
	}

	ip, _, _ := net.SplitHostPort(t.Server.Addr)
	t.Server.Addr = net.JoinHostPort(ip, port)

	return nil
}

func (t *DefaultServer) SetListenType(typ string) error {
	switch typ {
	case "http":
		t.ListenType = ListenTypeHttp
	case "https":
		t.ListenType = ListenTypeHttps
	default:
		return fmt.Errorf("invalid type - %s", typ)
	}

	return nil
}

func (t *DefaultServer) SetCertFile(filename string) {
	t.CertFile = filename
}

func (t *DefaultServer) SetKeyFile(filename string) {
	t.KeyFile = filename
}

func (t *DefaultServer) SetErrorLog(filename string) error {
	file, err := utils.OpenOrCreateFile(filename)
	if err != nil {
		return fmt.Errorf("can't open error log file - %s", filename)
	}

	t.ErrorLog = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

func (t *DefaultServer) SetDebugLog(filename string) error {
	file, err := utils.OpenOrCreateFile(filename)
	if err != nil {
		return fmt.Errorf("can't open debug log file - %s", filename)
	}

	t.DebugLog = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

func (t *DefaultServer) SetReadTimeout(timeout string) error {
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("invalid read timeout - %s", timeout)
	}

	t.Server.ReadTimeout = dur

	return nil
}

func (t *DefaultServer) SetReadHeaderTimeout(timeout string) error {
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("invalid read header timeout - %s", timeout)
	}

	t.Server.ReadHeaderTimeout = dur

	return nil
}

func (t *DefaultServer) SetWriteTimeout(timeout string) error {
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("invalid write timeout - %s", timeout)
	}

	t.Server.WriteTimeout = dur

	return nil
}

func (t *DefaultServer) SetIdleTimeout(timeout string) error {
	dur, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("invalid idle timeout - %s", timeout)
	}

	t.Server.IdleTimeout = dur

	return nil
}

func (t *DefaultServer) LoadConfig(filename string) error {
	return conf.DefaultDecoder.Decode(t.Config, filename)
}

func (t *DefaultServer) ListenAndServe() error {
	switch t.ListenType {
	case ListenTypeHttp:
		return t.Server.ListenAndServe()
	case ListenTypeHttps:
		return t.Server.ListenAndServeTLS(t.CertFile, t.KeyFile)
	default:
		return fmt.Errorf("unavailable listen type - %v", t.ListenType)
	}
}

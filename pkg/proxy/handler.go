package proxy

import (
	"errors"
	"github.com/abbot/go-http-auth"
	"net"
	"os"
)

type authType int

const (
	AuthTypeNone authType = iota
	AuthTypeBasic
)

func AuthTypeFromString(authTypeStr *string) (authType, error) {
	if authTypeStr == nil {
		return AuthTypeNone, nil
	}

	switch *authTypeStr {
	case "basic":
		return AuthTypeBasic, nil
	default:
		return -1, errors.New("unavailable auth type")
	}
}

type Handler struct {
	AccessLogger         Logger
	ErrorLogger          Logger
	AuthType             authType
	Htpasswd             auth.BasicAuth
	HtpasswdForRedirects auth.BasicAuth
	UseIpV4              *net.IP
	UseIpV6              *net.IP
	EnableUseIpHeader    bool
	BlockRequests        bool
	RedirectToProxy      *Proxy
	Conditions           map[Condition]Handler

	accessLoggerFile *os.File
	errorLoggerFile  *os.File
}

func (t *Handler) Close() error {
	if err := t.AccessLogger.Close(); err != nil {
		return err
	}
	if err := t.ErrorLogger.Close(); err != nil {
		return err
	}

	return nil
}

func (t *Handler) setFromConfig(config ConfigHandler) error {
	var err error

	if err = t.AccessLogger.Close(); err != nil {
		return err
	}
	if config.AccessLog != nil {
		if err = t.AccessLogger.SetFile(*config.AccessLog); err != nil {
			return err
		}
	}

	if err = t.ErrorLogger.Close(); err != nil {
		return err
	}
	if config.ErrorLog != nil {
		if err = t.ErrorLogger.SetFile(*config.ErrorLog); err != nil {
			return err
		}
	}

	if t.AuthType, err = AuthTypeFromString(config.AuthType); err != nil {
		return err
	}

	return nil
}

package proxy

import (
	"errors"
	"net"
	"os"
)

type authType int

const (
	AuthTypeNone authType = iota
	AuthTypeBasic
)

func AuthTypeFromString(authTypeStr string) (authType, error) {
	switch authTypeStr {
	case "none":
		return AuthTypeNone, nil
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
	Htpasswd             *BasicAuth
	HtpasswdForRedirects *BasicAuth
	UseIpV4              net.IP
	UseIpV6              net.IP
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
	if config.AuthType != nil {
		if t.AuthType, err = AuthTypeFromString(*config.AuthType); err != nil {
			return err
		}
	}
	if config.Htpasswd != nil {
		if t.Htpasswd, err = NewBasicAuth(*config.Htpasswd); err != nil {
			return err
		}
	}
	if config.HtpasswdForRedirects != nil {
		if t.HtpasswdForRedirects, err = NewBasicAuth(*config.HtpasswdForRedirects); err != nil {
			return err
		}
	}
	if config.UseIpV4 != nil {
		if t.UseIpV4 = net.ParseIP(*config.UseIpV4); t.UseIpV4 == nil {
			return errors.New("incorrect ipV4 address")
		}
	}
	if config.UseIpV6 != nil {
		if t.UseIpV6 = net.ParseIP(*config.UseIpV6); t.UseIpV6 == nil {
			return errors.New("incorrect ipV6 address")
		}
	}
	if config.EnableUseIpHeader != nil {
		t.EnableUseIpHeader = *config.EnableUseIpHeader
	}
	if config.BlockRequests != nil {
		t.BlockRequests = *config.BlockRequests
	}
	if config.RedirectToProxy != nil {
		t.RedirectToProxy = new(Proxy)
		t.RedirectToProxy.htpasswdForRedirects = t.HtpasswdForRedirects
		if err = t.RedirectToProxy.setFromConfig(*config.RedirectToProxy); err != nil {
			return err
		}
	}
	if config.Conditions != nil {
		t.Conditions = make(map[Condition]Handler)
		var configCondition ConfigCondition
		for configCondition = range config.Conditions {
			condition, err := NewCondition(configCondition.Condition)
			if err != nil {
				return err
			}

			handler := *t
			err = handler.setFromConfig(configCondition.Handler)
			if err != nil {
				return err
			}

			t.Conditions[*condition] = handler
		}
	}

	return nil
}

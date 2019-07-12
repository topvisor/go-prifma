package proxy

import (
	"errors"
	"net"
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
	AuthType             *authType
	Htpasswd             *BasicAuth
	HtpasswdForRedirects *BasicAuth
	UseIpV4              net.IP
	UseIpV6              net.IP
	EnableUseIpHeader    *bool
	BlockRequests        *bool
	RedirectToProxy      *Proxy

	conditions map[Condition]Handler
}

func (t *Handler) Close() error {
	if err := t.AccessLogger.Close(); err != nil {
		return err
	}
	if err := t.ErrorLogger.Close(); err != nil {
		return err
	}

	for _, handler := range t.conditions {
		if err := handler.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (t *Handler) SetConditionHandler(cond Condition, handler Handler) error {
	if t.conditions == nil {
		t.conditions = make(map[Condition]Handler)
	}

	if oldHandler, exists := t.conditions[cond]; exists {
		if err := oldHandler.Close(); err != nil {
			return err
		}
	}

	if !handler.AccessLogger.IsInited() {
		handler.AccessLogger = t.AccessLogger
	}
	if !handler.ErrorLogger.IsInited() {
		handler.ErrorLogger = t.ErrorLogger
	}
	if handler.AuthType == nil {
		handler.AuthType = t.AuthType
	}
	if handler.Htpasswd == nil {
		handler.Htpasswd = t.Htpasswd
	}
	if handler.HtpasswdForRedirects == nil {
		handler.HtpasswdForRedirects = t.HtpasswdForRedirects
	}
	if handler.UseIpV4 == nil {
		handler.UseIpV4 = t.UseIpV4
	}
	if handler.UseIpV6 == nil {
		handler.UseIpV6 = t.UseIpV6
	}
	if handler.EnableUseIpHeader == nil {
		handler.EnableUseIpHeader = t.EnableUseIpHeader
	}
	if handler.BlockRequests == nil {
		handler.BlockRequests = t.BlockRequests
	}
	if handler.RedirectToProxy == nil {
		handler.RedirectToProxy = t.RedirectToProxy
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
		authType, err := AuthTypeFromString(*config.AuthType)
		if err != nil {
			return err
		}

		t.AuthType = &authType
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

	t.EnableUseIpHeader = config.EnableUseIpHeader
	t.BlockRequests = config.BlockRequests

	if config.RedirectToProxy != nil {
		t.RedirectToProxy = new(Proxy)
		t.RedirectToProxy.htpasswdForRedirects = t.HtpasswdForRedirects
		if err = t.RedirectToProxy.setFromConfig(*config.RedirectToProxy); err != nil {
			return err
		}
	}
	if config.Conditions != nil {
		for _, configCondition := range config.Conditions {
			condition, err := NewCondition(configCondition.Condition)
			if err != nil {
				return err
			}

			handler := Handler{}
			if err = handler.setFromConfig(configCondition.Handler); err != nil {
				return err
			}

			if err = t.SetConditionHandler(*condition, handler); err != nil {
				return nil
			}
		}
	}

	return nil
}

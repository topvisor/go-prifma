package proxy

import (
	"errors"
	"net"
)

type Handler struct {
	AccessLogger      Logger
	ErrorLogger       Logger
	Htpasswd          *BasicAuth
	EnableBasicAuth   *bool
	OutgoingIpV4      net.IP
	OutgoingIpV6      net.IP
	EnableUseIpHeader *bool
	BlockRequests     *bool
	Proxy             *Proxy

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
	if config.Htpasswd != nil {
		if t.Htpasswd, err = NewBasicAuth(*config.Htpasswd); err != nil {
			return err
		}
	}

	t.EnableBasicAuth = config.EnableBasicAuth

	if config.OutgoingIpV4 != nil {
		if t.OutgoingIpV4 = net.ParseIP(*config.OutgoingIpV4); t.OutgoingIpV4 == nil {
			return errors.New("incorrect ipV4 address")
		}
	}
	if config.OutgoingIpV6 != nil {
		if t.OutgoingIpV6 = net.ParseIP(*config.OutgoingIpV6); t.OutgoingIpV6 == nil {
			return errors.New("incorrect ipV6 address")
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

			if err = t.SetConditionHandler(*condition, handler); err != nil {
				return nil
			}
		}
	}

	return nil
}

package proxy

import (
	"encoding/json"
	"os"
)

// ConfigProxy is a part of config.json which describes a Proxy
type ConfigProxy struct {
	Url          string            `json:"url"`
	ProxyHeaders map[string]string `json:"proxyHeaders,omitempty"`
}

// ConfigCondition is a part of config.json which describes a Condition
type ConfigCondition struct {
	Condition string        `json:"condition"`
	Handler   ConfigHandler `json:"handler"`
}

// ConfigListen is a part of config.json which describes a Server
type ConfigListen struct {
	ListenIp   *string `json:"listenIp,omitempty"`
	ListenPort int     `json:"listenPort"`
	ListenType string  `json:"listenType"`
}

// ConfigListen is a part of config.json which describes a Handler
type ConfigHandler struct {
	AccessLog         *string           `json:"accessLog,omitempty"`
	ErrorLog          *string           `json:"errorLog,omitempty"`
	DialTimeout       *int              `json:"dialTimeout,omitempty"`
	Htpasswd          *string           `json:"htpasswd,omitempty"`
	EnableBasicAuth   *bool             `json:"enableBasicAuth,omitempty"`
	OutgoingIpV4      *string           `json:"outgoingIpV4,omitempty"`
	OutgoingIpV6      *string           `json:"outgoingIpV6,omitempty"`
	EnableUseIpHeader *bool             `json:"enableUseIpHeader,omitempty"`
	BlockRequests     *bool             `json:"blockRequests,omitempty"`
	Proxy             *ConfigProxy      `json:"proxy,omitempty"`
	Conditions        []ConfigCondition `json:"conditions,omitempty"`
}

// ConfigListen is a part of config.json which describes the Server and the base Handler
type Config struct {
	Listen ConfigListen `json:"server"`
	ConfigHandler
}

func ParseConfig(jsonStr string) (*Config, error) {
	config := new(Config)
	if err := json.Unmarshal([]byte(jsonStr), config); err != nil {
		return nil, err
	}

	return config, nil
}

func ParseConfigFromFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return config, nil
}

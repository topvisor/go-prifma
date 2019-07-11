package proxy

import (
	"encoding/json"
	"os"
)

type ConfigProxy struct {
	Url       string            `json:"url"`
	BasicAuth *string           `json:"basicAuth,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

type ConfigCondition struct {
	Condition string        `json:"condition"`
	Handler   ConfigHandler `json:"handler"`
}

type ConfigServer struct {
	ListenPort int    `json:"listenPort"`
	ListenType string `json:"listenType"`
}

type ConfigHandler struct {
	AccessLog            *string            `json:"accessLog,omitempty"`
	ErrorLog             *string            `json:"errorLog,omitempty"`
	AuthType             *string            `json:"authType,omitempty"`
	Htpasswd             *string            `json:"htpasswd,omitempty"`
	HtpasswdForRedirects *string            `json:"htpasswdForRedirects,omitempty"`
	UseIpV4              *string            `json:"useIpV4,omitempty"`
	UseIpV6              *string            `json:"useIpV6,omitempty"`
	EnableUseIpHeader    bool               `json:"enableUseIpHeader,omitempty"`
	BlockRequests        bool               `json:"blockRequests,omitempty"`
	RedirectToProxy      *ConfigProxy       `json:"redirectToProxy,omitempty"`
	Conditions           *[]ConfigCondition `json:"conditions,omitempty"`
}

type Config struct {
	Server ConfigServer `json:"server"`
	ConfigHandler
}

func ParseConfig(jsonStr string) (*Config, error) {
	config := Config{}
	if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func ParseConfigFromFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	config := Config{}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

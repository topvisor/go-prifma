package proxyserver

import (
	"encoding/json"
	"os"
)

type ConfigProxy struct {
	Url       string             `json:"url"`
	BasicAuth *string            `json:"basicAuth,omitempty"`
	Headers   *map[string]string `json:"headers,omitempty"`
}

type ConfigCondition struct {
	Condition string `json:"condition"`
	Config    Config `json:"config"`
}

type Config struct {
	AccessLog         *string            `json:"accessLog,omitempty"`
	ErrorLog          *string            `json:"errorLog,omitempty"`
	ListenPort        *int               `json:"listenPort,omitempty"`
	ListenType        *string            `json:"listenType,omitempty"`
	AuthType          *string            `json:"authType,omitempty"`
	Users             *string            `json:"users,omitempty"`
	UsersForRedirects *string            `json:"usersForRedirects,omitempty"`
	UseIpV4           *string            `json:"useIpV4,omitempty"`
	UseIpV6           *string            `json:"useIpV6,omitempty"`
	EnableUseIpHeader *bool              `json:"enableUseIpHeader,omitempty"`
	BlockRequests     *bool              `json:"blockRequests,omitempty"`
	RedirectToProxy   *ConfigProxy       `json:"redirectToProxy,omitempty"`
	Conditions        *[]ConfigCondition `json:"conditions,omitempty"`
}

func NewConfig(jsonStr *string) (*Config, error) {
	var config *Config
	if err := json.Unmarshal([]byte(*jsonStr), config); err != nil {
		return nil, err
	}

	return config, nil
}

func NewConfigFromFile(filename *string) (*Config, error) {
	file, err := os.Open(*filename)
	if err != nil {
		return nil, err
	}

	var config *Config
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

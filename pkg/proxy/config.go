package proxy

import (
	"encoding/json"
	"os"
	"reflect"
)

type ConfigOutgoingIp struct {
	Ips []string
}

func (t *ConfigOutgoingIp) UnmarshalJSON(dataBytes []byte) error {
	var data interface{}

	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return err
	}

	switch data.(type) {
	case string:
		t.Ips = []string{data.(string)}
	case []interface{}:
		t.Ips = make([]string, len(data.([]interface{})))
		for i, ip := range data.([]interface{}) {
			if ipStr, ok := ip.(string); ok {
				t.Ips[i] = ipStr
			} else {
				return &json.InvalidUnmarshalError{Type: reflect.TypeOf(t)}
			}
		}
	default:
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(t)}
	}

	return nil
}

func (t *ConfigOutgoingIp) MarshalJSON() ([]byte, error) {
	var data interface{}

	switch len(t.Ips) {
	case 0:
		data = nil
	case 1:
		data = t.Ips[0]
	default:
		data = t.Ips
	}

	return json.Marshal(data)
}

// ConfigProxy is a part of config.json which describes a Proxy
type ConfigProxy struct {
	Url          string            `json:"url"`
	ProxyHeaders map[string]string `json:"proxyHeaders"`
}

// ConfigCondition is a part of config.json which describes a Condition
type ConfigCondition struct {
	Condition string        `json:"condition"`
	Handler   ConfigHandler `json:"handler"`
}

// ConfigListen is a part of config.json which describes a server
type ConfigListen struct {
	ListenIp          *string `json:"listenIp"`
	ListenPort        int     `json:"listenPort"`
	ListenType        string  `json:"listenType"`
	ErrorLog          *string `json:"errorLog"`
	ReadTimeout       *string `json:"readTimeout"`
	ReadHeaderTimeout *string `json:"readHeaderTimeout"`
	WriteTimeout      *string `json:"writeTimeout"`
	IdleTimeout       *string `json:"idleTimeout"`
}

// ConfigListen is a part of config.json which describes a Handler
type ConfigHandler struct {
	AccessLog         *string           `json:"accessLog"`
	DumpLog           *string           `json:"dumpLog"`
	Htpasswd          *string           `json:"htpasswd"`
	EnableBasicAuth   *bool             `json:"enableBasicAuth"`
	OutgoingIpV4      *ConfigOutgoingIp `json:"outgoingIpV4"`
	OutgoingIpV6      *ConfigOutgoingIp `json:"outgoingIpV6"`
	EnableUseIpHeader *bool             `json:"enableUseIpHeader"`
	BlockRequests     *bool             `json:"blockRequests"`
	Proxy             *ConfigProxy      `json:"proxy"`
	Conditions        []ConfigCondition `json:"conditions"`
}

// ConfigListen is a part of config.json which describes the server and the base Handler
type Config struct {
	Listen ConfigListen `json:"server"`
	ConfigHandler
}

// ParseConfig parses a config from the json string
func ParseConfig(jsonStr string) (*Config, error) {
	config := new(Config)
	if err := json.Unmarshal([]byte(jsonStr), config); err != nil {
		return nil, err
	}

	return config, nil
}

// ParseConfigFromFile parses a config from the file
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

package prifma

import (
	"encoding/json"
	"github.com/topvisor/prifma/pkg/conf"
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
	Url          string
	ProxyHeaders map[string]string
}

func (t *ConfigProxy) Call(name string, args ...string) error {
	if name != "header" || len(args) != 2 {
		return NewErrWrongCall(name, args)
	}
	if t.ProxyHeaders == nil {
		t.ProxyHeaders = make(map[string]string)
	}

	t.ProxyHeaders[args[0]] = t.ProxyHeaders[args[1]]

	return nil
}

func (t *ConfigProxy) CallBlock(name string, args ...string) (conf.Block, error) {
	return nil, NewErrWrongCall(name, args)
}

// ConfigCondition is a part of config.json which describes a Condition
type ConfigCondition struct {
	Key   string
	Type  string
	Value string
}

// ConfigListen is a part of config.json which describes a prifma
type ConfigListen struct {
	ListenIp          string
	ListenPort        string
	ListenType        string
	ErrorLog          string
	ReadTimeout       string
	ReadHeaderTimeout string
	WriteTimeout      string
	IdleTimeout       string
}

func (t *ConfigListen) Call(name string, args ...string) error {
	if len(args) != 1 {
		return NewErrWrongCall(name, args)
	}

	switch name {
	case "listen_ip":
		t.ListenIp = args[0]
	case "listen_port":
		t.ListenPort = args[0]
	case "listen_type":
		t.ListenType = args[0]
	case "error_log":
		t.ErrorLog = args[0]
	case "read_timeout":
		t.ReadTimeout = args[0]
	case "read_header_timeout":
		t.ReadHeaderTimeout = args[0]
	case "write_timeout":
		t.WriteTimeout = args[0]
	case "idle_timeout":
		t.IdleTimeout = args[0]
	default:
		return NewErrWrongCall(name, args)
	}

	return nil
}

func (t *ConfigListen) CallBlock(name string, args ...string) (conf.Block, error) {
	return nil, NewErrWrongCall(name, args)
}

// ConfigListen is a part of config.json which describes a Handler
type ConfigHandler struct {
	Condition           *ConfigCondition
	AccessLog           string
	DumpLog             string
	BasicAuth           string
	OutgoingIpV4        []string
	OutgoingIpV6        []string
	EnableUseIpHeader   bool
	EnableBlockRequests bool
	Proxy               *ConfigProxy
	Handlers            []*ConfigHandler
}

func (t *ConfigHandler) Call(name string, args ...string) error {
	switch name {
	case "outgoing_ip_v4":
		t.OutgoingIpV4 = args
	case "outgoing_ip_v6":
		t.OutgoingIpV6 = args
	default:
		switch len(args) {
		case 0:
			switch name {
			case "enable_use_ip_header":
				t.EnableUseIpHeader = true
			case "disable_use_ip_header":
				t.EnableUseIpHeader = false
			case "enable_block_requests":
				t.EnableBlockRequests = true
			case "disable_block_requests":
				t.EnableBlockRequests = false
			case "disable_proxy_requests":
				t.Proxy = nil
			default:
				return NewErrWrongCall(name, args)
			}
		case 1:
			switch name {
			case "access_log":
				t.AccessLog = args[0]
			case "dump_log":
				t.DumpLog = args[0]
			case "basic_auth":
				t.BasicAuth = args[0]
			case "enable_proxy_requests":
				t.Proxy = &ConfigProxy{
					Url: args[0],
				}
			default:
				return NewErrWrongCall(name, args)
			}
		default:
			return NewErrWrongCall(name, args)
		}
	}

	return nil
}

func (t *ConfigHandler) CallBlock(name string, args ...string) (conf.Block, error) {
	switch name {
	case "condition":
		if len(args) != 3 {
			return nil, NewErrWrongCall(name, args)
		}

		handler := *t
		handler.Condition = &ConfigCondition{
			Key:   args[0],
			Type:  args[1],
			Value: args[2],
		}

		if t.Handlers == nil {
			t.Handlers = []*ConfigHandler{&handler}
		} else {
			t.Handlers = append(t.Handlers, &handler)
		}

		return &handler, nil
	case "enable_proxy_requests":
		if len(args) != 1 {
			return nil, NewErrWrongCall(name, args)
		}

		t.Proxy = &ConfigProxy{
			Url: args[0],
		}

		return t.Proxy, nil
	default:
		return nil, NewErrWrongCall(name, args)
	}
}

// ConfigListen is a part of config.json which describes the prifma and the base Handler
type Config struct {
	Listen *ConfigListen
	ConfigHandler
}

// ParseConfig parses a config from the file
func ParseConfig(filename string) (*Config, error) {
	t := new(Config)

	return t, conf.DefaultDecoder.Decode(t, filename)
}

func (t *Config) CallBlock(name string, args ...string) (conf.Block, error) {
	if name != "prifma" || len(args) != 0 || t.Listen != nil {
		return nil, NewErrWrongCall(name, args)
	}

	t.Listen = new(ConfigListen)

	return t.Listen, nil
}

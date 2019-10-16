package prifma_new

import (
	"github.com/topvisor/prifma/pkg/conf"
)

type ConfigMain struct {
	ConfigServer conf.Block
	ConfigModule conf.Block
}

func NewConfigMain(server Server) conf.Block {
	return &ConfigMain{
		ConfigServer: NewConfigServer(server),
		ConfigModule: NewConfigModule(server),
	}
}

func (t *ConfigMain) Call(command conf.Command) error {
	return t.ConfigModule.Call(command)
}

func (t *ConfigMain) CallBlock(command conf.Command) (conf.Block, error) {
	switch command.GetName() {
	case "server":
		return t.ConfigServer.CallBlock(command)
	default:
		return t.ConfigModule.CallBlock(command)
	}
}

type ConfigServer struct {
	Server Server
}

func NewConfigServer(server Server) conf.Block {
	return &ConfigServer{
		Server: server,
	}
}

func (t *ConfigServer) Call(command conf.Command) error {
	if len(command.GetArgs()) != 1 {
		return NewErrWrongDirective(command)
	}

	arg := command.GetArgs()[0]

	switch command.GetName() {
	case "listen_ip":
		return t.Server.SetListenIp(arg)
	case "listen_port":
		return t.Server.SetListenPort(arg)
	case "listen_schema":
		return t.Server.SetListenType(arg)
	case "error_log":
		return t.Server.SetErrorLog(arg)
	case "read_timeout":
		return t.Server.SetReadTimeout(arg)
	case "read_header_timeout":
		return t.Server.SetReadHeaderTimeout(arg)
	case "write_timeout":
		return t.Server.SetWriteTimeout(arg)
	case "idle_timeout":
		return t.Server.SetIdleTimeout(arg)
	}

	return NewErrWrongDirective(command)
}

func (t *ConfigServer) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() != "server" || len(command.GetArgs()) != 0 {
		return nil, NewErrWrongDirective(command)
	}

	return t, nil
}

type ConfigModule struct {
	Server Server
}

func NewConfigModule(server Server) conf.Block {
	return &ConfigModule{
		Server: server,
	}
}

func (t *ConfigModule) Call(command conf.Command) error {
	for _, module := range t.Server.GetModules() {
		if err := module.Call(command); err != nil {
			if _, ok := err.(*ErrModuleDirectiveNotFound); !ok {
				return err
			}
		} else {
			return nil
		}
	}

	return NewErrWrongDirective(command)
}

func (t *ConfigModule) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() == "condition" {
		if len(command.GetArgs()) != 3 {
			return nil, NewErrWrongDirective(command)
		}

	}

	for _, module := range t.Server.GetModules() {
		if block, err := module.CallBlock(command); err != nil {
			if _, ok := err.(*ErrModuleDirectiveNotFound); !ok {
				return nil, err
			}
		} else {
			return block, nil
		}
	}

	return nil, NewErrWrongDirective(command)
}

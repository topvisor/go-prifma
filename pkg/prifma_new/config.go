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
		return conf.NewCommandError(command)
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

	return conf.NewCommandError(command)
}

func (t *ConfigServer) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() != "server" || len(command.GetArgs()) != 0 {
		return nil, conf.NewCommandError(command)
	}

	return t, nil
}

type ConfigModule struct {
	Server Server
	Conds  []Condition
}

func NewConfigModule(server Server) conf.Block {
	return &ConfigModule{
		Server: server,
		Conds:  make([]Condition, 0),
	}
}

func (t *ConfigModule) Call(command conf.Command) error {
	module := t.Server.GetModulesManager().GetModules(t.Conds...)[command.GetName()]
	if module == nil {
		return conf.NewCommandError(command)
	}

	return module.Call(command)
}

func (t *ConfigModule) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() == "condition" {
		return t.CallCondition(command)
	}

	module := t.Server.GetModulesManager().GetModules(t.Conds...)[command.GetName()]
	if module == nil {
		return nil, conf.NewCommandError(command)
	}

	return module.CallBlock(command)
}

func (t *ConfigModule) CallCondition(command conf.Command) (conf.Block, error) {
	args := command.GetArgs()
	if len(args) != 3 {
		return nil, conf.NewCommandError(command)
	}

	cond, err := NewCondition(args[0], args[1], args[2])
	if err != nil {
		return nil, err
	}

	conditionBlock := &ConfigModule{
		Server: t.Server,
		Conds:  append(t.Conds, cond),
	}

	return conditionBlock, nil
}

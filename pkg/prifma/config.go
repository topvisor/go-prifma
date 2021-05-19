package prifma

import (
	"github.com/topvisor/go-prifma/pkg/conf"
)

type ConfigMain struct {
	ConfigServer conf.Block
	ConfigModule conf.Block
}

func NewConfigMain(server Server) *ConfigMain {
	return &ConfigMain{
		ConfigServer: NewConfigServer(server),
		ConfigModule: NewConfigModule(server.GetModulesManager()),
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

func NewConfigServer(server Server) *ConfigServer {
	return &ConfigServer{
		Server: server,
	}
}

func (t *ConfigServer) Call(command conf.Command) (err error) {
	if len(command.GetArgs()) != 1 {
		return conf.NewErrCommandArgsNumber(command)
	}

	arg := command.GetArgs()[0]

	switch command.GetName() {
	case "listen_ip":
		err = t.Server.SetListenIp(arg)
	case "listen_port":
		err = t.Server.SetListenPort(arg)
	case "listen_schema":
		err = t.Server.SetListenType(arg)
	case "cert_file":
		t.Server.SetCertFile(arg)
	case "key_file":
		t.Server.SetKeyFile(arg)
	case "error_log":
		err = t.Server.SetErrorLog(arg)
	case "debug_log":
		err = t.Server.SetDebugLog(arg)
	case "read_timeout":
		err = t.Server.SetReadTimeout(arg)
	case "read_header_timeout":
		err = t.Server.SetReadHeaderTimeout(arg)
	case "write_timeout":
		err = t.Server.SetWriteTimeout(arg)
	case "idle_timeout":
		err = t.Server.SetIdleTimeout(arg)
	default:
		return conf.NewErrCommandName(command)
	}

	if err != nil {
		err = conf.NewErrCommand(command, err.Error())
	}

	return err
}

func (t *ConfigServer) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() != "server" {
		return nil, conf.NewErrCommandName(command)
	}

	if len(command.GetArgs()) != 0 {
		return nil, conf.NewErrCommandArgsNumber(command)
	}

	return t, nil
}

type ConfigModule struct {
	ModulesManager ModulesManager
	Conds          []Condition
}

func NewConfigModule(modulesManager ModulesManager) *ConfigModule {
	return &ConfigModule{
		ModulesManager: modulesManager,
		Conds:          make([]Condition, 0),
	}
}

func (t *ConfigModule) Call(command conf.Command) error {
	module := t.ModulesManager.GetModule(command.GetName(), t.Conds...)
	if module == nil {
		return conf.NewErrCommandName(command)
	}

	return module.Call(command)
}

func (t *ConfigModule) CallBlock(command conf.Command) (conf.Block, error) {
	if command.GetName() == "condition" {
		return t.CallCondition(command)
	}

	module := t.ModulesManager.GetModule(command.GetName(), t.Conds...)
	if module == nil {
		return nil, conf.NewErrCommandName(command)
	}

	return module.CallBlock(command)
}

func (t *ConfigModule) CallCondition(command conf.Command) (conf.Block, error) {
	args := command.GetArgs()
	if len(args) != 3 {
		return nil, conf.NewErrCommandArgsNumber(command)
	}

	cond, err := NewCondition(args[0], args[1], args[2])
	if err != nil {
		return nil, conf.NewErrCommand(command, err.Error())
	}

	conditionBlock := &ConfigModule{
		ModulesManager: t.ModulesManager,
		Conds:          append(t.Conds, cond),
	}

	return conditionBlock, nil
}

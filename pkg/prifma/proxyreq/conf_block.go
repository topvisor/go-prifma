package proxyreq

import (
	"github.com/topvisor/go-prifma/pkg/conf"
	"net/http"
)

const ModuleBlockDirective = "proxy_header"

type ConfBlock struct {
	Header *http.Header
}

func NewConfBlock(header *http.Header) conf.Block {
	return &ConfBlock{
		Header: header,
	}
}

func (t *ConfBlock) Call(command conf.Command) error {
	args := command.GetArgs()

	if command.GetName() != ModuleBlockDirective || len(args) != 2 {
		return conf.NewErrCommand(command)
	}

	key := args[0]
	val := args[1]

	if val == "" {
		t.Header.Del(key)
	} else {
		t.Header.Set(key, val)
	}

	return nil
}

func (t *ConfBlock) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}

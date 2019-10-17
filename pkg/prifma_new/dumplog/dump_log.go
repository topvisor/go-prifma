package dumplog

import (
	"fmt"
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma_new"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

const ModuleDirective = "dump_log"

type DumpLog struct {
	Logger *log.Logger
}

func New() prifma_new.Module {
	return new(DumpLog)
}

func (t *DumpLog) BeforeHandleRequest(req *http.Request) error {
	if t.Logger == nil {
		return nil
	}

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}

	t.Logger.Println(strings.TrimSpace(string(dump)))

	return nil
}

func (t *DumpLog) Off() error {
	t.Logger = nil

	return nil
}

func (t *DumpLog) SetFilename(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("can't open dump log file: '%s'", filename)
	}

	t.Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

func (t *DumpLog) GetDirective() string {
	return ModuleDirective
}

func (t *DumpLog) Clone() prifma_new.Module {
	clone := *t

	return &clone
}

func (t *DumpLog) Call(command conf.Command) error {
	if command.GetName() != ModuleDirective || len(command.GetArgs()) != 1 {
		return conf.NewErrCommand(command)
	}

	arg := command.GetArgs()[0]
	if arg == "off" {
		return t.Off()
	}

	return t.SetFilename(arg)
}

func (t *DumpLog) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}

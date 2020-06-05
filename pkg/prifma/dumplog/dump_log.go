package dumplog

import (
	"fmt"
	"github.com/topvisor/go-prifma/pkg/conf"
	"github.com/topvisor/go-prifma/pkg/prifma"
	"github.com/topvisor/go-prifma/pkg/utils"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

const ModuleDirective = "dump_log"

type DumpLog struct {
	Logger *log.Logger
}

func New() *DumpLog {
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
	file, err := utils.OpenOrCreateFile(filename)
	if err != nil {
		return fmt.Errorf("can't open dump log file: '%s'", filename)
	}

	t.Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

func (t *DumpLog) GetDirective() string {
	return ModuleDirective
}

func (t *DumpLog) Clone() prifma.Module {
	clone := *t

	return &clone
}

func (t *DumpLog) Call(command conf.Command) error {
	if command.GetName() != ModuleDirective {
		return conf.NewErrCommandName(command)
	}

	if len(command.GetArgs()) != 1 {
		return conf.NewErrCommandArgsNumber(command)
	}

	arg := command.GetArgs()[0]
	if arg == "off" {
		return t.Off()
	}

	return t.SetFilename(arg)
}

func (t *DumpLog) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommandMustHaveNoBlock(command)
}

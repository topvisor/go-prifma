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

func (t *DumpLog) Off() {
	t.Logger = nil
}

func (t *DumpLog) SetFilename(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("can't open dump log file: '%s'", filename)
	}

	t.Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

func (t *DumpLog) Call(command conf.Command) error {
	if command.GetName() != "dump_log" {
		return prifma_new.NewErrModuleDirectiveNotFound(command)
	}
	if len(command.GetArgs()) != 1 {
		return prifma_new.NewErrWrongDirective(command)
	}

	arg := command.GetArgs()[0]

	if arg == "off" {
		t.Off()

		return nil
	}

	return t.SetFilename(arg)
}

func (t *DumpLog) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, prifma_new.NewErrModuleDirectiveNotFound(command)
}
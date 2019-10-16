package accesslog

import (
	"fmt"
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma_new"
	"log"
	"net/http"
	"os"
)

type AccessLog struct {
	Logger *log.Logger
}

func New() prifma_new.Module {
	return new(AccessLog)
}

func (t *AccessLog) Off() {
	t.Logger = nil
}

func (t *AccessLog) SetFilename(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("can't open access log file: '%s'", filename)
	}

	t.Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

func (t *AccessLog) AfterWriteResponse(req *http.Request, resp prifma_new.Response) error {
	if t.Logger == nil {
		return nil
	}

	var user *string
	if username, _, ok := prifma_new.ProxyBasicAuth(req); ok {
		user = &username
	}

	t.Logger.Printf(
		"%s %d %s %s %v l/%v r/%v\n",
		req.RemoteAddr,
		resp.GetCode(),
		req.Method,
		req.RequestURI,
		user,
		resp.GetLAddr(),
		resp.GetRAddr(),
	)

	return nil
}

func (t *AccessLog) Call(command conf.Command) error {
	if command.GetName() != "access_log" {
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

func (t *AccessLog) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, prifma_new.NewErrModuleDirectiveNotFound(command)
}

package basicauth

import (
	"encoding/csv"
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"github.com/topvisor/prifma/pkg/conf"
	"github.com/topvisor/prifma/pkg/prifma_new"
	"github.com/topvisor/prifma/pkg/utils"
	"os"
)

const ModuleDirective = "basic_auth"

type BasicAuth struct {
	Users map[string]string
}

func New() prifma_new.Module {
	return new(BasicAuth)
}

func (t *BasicAuth) HandleRequest(result prifma_new.HandleRequestResult) (prifma_new.HandleRequestResult, error) {
	if t.Users == nil {
		return result, nil
	}

	user, pass, _ := utils.ProxyBasicAuth(result.GetRequest())
	secret, ok := t.Users[user]
	if !ok || !auth.CheckSecret(pass, secret) {
		result.SetResponse(NewResponseRequireAuth(result.GetRequest()))
	}

	return result, nil
}

func (t *BasicAuth) Off() error {
	t.Users = nil

	return nil
}

func (t *BasicAuth) LoadHtpasswdFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("can't open htpasswd file: '%s'", filename)
	}

	reader := csv.NewReader(file)
	reader.Comma = ':'
	reader.Comment = '#'
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("wrong format of htpasswd file: '%s'", filename)
	}

	t.Users = make(map[string]string)
	for _, record := range records {
		t.Users[record[0]] = record[1]
	}

	return nil
}

func (t *BasicAuth) GetDirective() string {
	return ModuleDirective
}

func (t *BasicAuth) Clone() prifma_new.Module {
	clone := *t

	return &clone
}

func (t *BasicAuth) Call(command conf.Command) error {
	if command.GetName() != ModuleDirective || len(command.GetArgs()) != 1 {
		return conf.NewErrCommand(command)
	}

	arg := command.GetArgs()[0]
	if arg == "off" {
		return t.Off()
	}

	return t.LoadHtpasswdFile(arg)
}

func (t *BasicAuth) CallBlock(command conf.Command) (conf.Block, error) {
	return nil, conf.NewErrCommand(command)
}

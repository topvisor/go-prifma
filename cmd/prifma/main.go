package main

import (
	"github.com/topvisor/prifma/pkg/prifma_new"
	"github.com/topvisor/prifma/pkg/prifma_new/accesslog"
	"github.com/topvisor/prifma/pkg/prifma_new/basicauth"
	"github.com/topvisor/prifma/pkg/prifma_new/dumplog"
)

func main() {
	flags, err := parseFlags()
	if err != nil {
		panic(err)
	}

	if flags.help {
		flags.PrintDefaults()
	}

	if err = start(flags.config); err != nil {
		panic(err)
	}
}

func start(configFilename string) error {
	server := prifma_new.NewServer(
		dumplog.New(),
		basicauth.New(),
		accesslog.New(),
	)

	if err := server.LoadConfig(configFilename); err != nil {
		return err
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

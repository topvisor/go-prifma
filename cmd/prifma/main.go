package main

import (
	"github.com/topvisor/prifma/pkg/prifma"
	"github.com/topvisor/prifma/pkg/prifma/accesslog"
	"github.com/topvisor/prifma/pkg/prifma/basicauth"
	"github.com/topvisor/prifma/pkg/prifma/blockreq"
	"github.com/topvisor/prifma/pkg/prifma/dumplog"
	"github.com/topvisor/prifma/pkg/prifma/http"
	"github.com/topvisor/prifma/pkg/prifma/outgoingip"
	"github.com/topvisor/prifma/pkg/prifma/proxyreq"
	"github.com/topvisor/prifma/pkg/prifma/tunnel"
	"github.com/topvisor/prifma/pkg/prifma/useipheader"
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
	server := prifma.NewServer(
		dumplog.New(),
		blockreq.New(),
		basicauth.New(),
		outgoingip.New(),
		useipheader.New(),
		proxyreq.New(),
		accesslog.New(),
		tunnel.New(),
		http.New(),
	)

	if err := server.LoadConfig(configFilename); err != nil {
		return err
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

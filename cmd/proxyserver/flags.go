package main

import (
	"flag"
	"os"
)

const (
	configFlag      = "config"
	configFlagShort = "c"
	configDefault   = "./config.json"
	configUsage     = "Set path to config.json"

	initFlag    = "init"
	initDefault = false
	initUsage   = "Create default config.json (if doesn't exist). See all config.json params in README.md"

	helpFlag      = "help"
	helpFlagShort = "h"
	helpDefault   = false
	helpUsage     = "Show this help"

	listenFlag      = "listen"
	listenFlagShort = "l"
	listenDefault   = false
	listenUsage     = "Start listening and serving requests"
)

type flags struct {
	config string
	init   bool
	help   bool
	listen bool

	flag.FlagSet
}

func parseFlags() (*flags, error) {
	flags := new(flags)
	err := flags.parse()

	return flags, err
}

func (t *flags) parse() error {
	t.StringVar(&t.config, configFlag, configDefault, configUsage)
	t.StringVar(&t.config, configFlagShort, configDefault, shortUsage(configUsage))
	t.BoolVar(&t.init, initFlag, initDefault, initUsage)
	t.BoolVar(&t.help, helpFlag, helpDefault, helpUsage)
	t.BoolVar(&t.help, helpFlagShort, helpDefault, shortUsage(helpUsage))
	t.BoolVar(&t.listen, listenFlag, listenDefault, listenUsage)
	t.BoolVar(&t.listen, listenFlagShort, listenDefault, shortUsage(listenUsage))

	return t.Parse(os.Args[1:])
}

func shortUsage(usage string) string {
	return usage + " (shorthand)"
}

package main

import (
	"flag"
	"os"
)

const (
	configFlag      = "config"
	configFlagShort = "c"
	configDefault   = "./prifma.conf"
	configUsage     = "Set path to config.json"

	helpFlag      = "help"
	helpFlagShort = "h"
	helpDefault   = false
	helpUsage     = "Show this help"
)

type flags struct {
	config string
	help   bool

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
	t.BoolVar(&t.help, helpFlag, helpDefault, helpUsage)
	t.BoolVar(&t.help, helpFlagShort, helpDefault, shortUsage(helpUsage))

	return t.Parse(os.Args[1:])
}

func shortUsage(usage string) string {
	return usage + " (shorthand)"
}

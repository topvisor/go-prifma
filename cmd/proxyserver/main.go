package main

import (
	"encoding/json"
	"github.com/topvisor/go-proxy-server/pkg/proxy"
	"os"
)

func main() {
	flags, err := parseFlags()
	if err != nil {
		panic(err)
	}

	if flags.init {
		if err = createDefaultConfigJSON(flags.config); err != nil {
			panic(err)
		}
	}
	if flags.listen {
		if err = start(flags.config); err != nil {
			panic(err)
		}
	}
	if !flags.init && !flags.listen {
		flags.PrintDefaults()
	}
}

func start(configFilename string) error {
	server, err := proxy.NewServerFromConfigFile(configFilename)
	if err != nil {
		return err
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func createDefaultConfigJSON(filename string) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	config := proxy.Config{
		Listen: proxy.ConfigListen{
			ListenType: "http",
			ListenPort: 3128,
		},
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	if err = encoder.Encode(config); err != nil {
		return err
	}

	return file.Close()
}

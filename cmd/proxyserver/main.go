package main

import (
	"../../pkg/proxy"
	"encoding/json"
	"log"
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
	server := new(proxy.Server)
	if err := server.LoadFromConfig(configFilename); err != nil {
		return err
	}
	if !server.ErrorLogger.IsInited() {
		errorLogger := log.New(os.Stderr, "", log.LstdFlags)
		if err := server.ErrorLogger.SetLogger(errorLogger); err != nil {
			return err
		}
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

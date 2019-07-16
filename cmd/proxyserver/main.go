package main

import (
	"encoding/json"
	"fmt"
	"go-proxy-server/pkg/proxy"
	"os"
)

func main() {
	flags, err := parseFlags()
	if err != nil {
		panic(err)
	}

	if flags.help {
		flags.PrintDefaults()
	} else if flags.init {
		if err = createDefaultConfigJSON(flags.config); err != nil {
			panic(err)
		}
	} else if flags.listen {
		fmt.Println("ololo")
	} else {
		flags.PrintDefaults()
	}
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

	if err = json.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	return file.Close()
}

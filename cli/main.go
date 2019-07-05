package main

import (
	"../../proxy"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	filename := "test.json"
	testJson := map[string]interface{}{
		"listenPort": 31299,
		"listenType": "http",
	}

	jsonData, err := json.Marshal(testJson)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filename, jsonData, 0664)
	if err != nil {
		log.Fatal(err)
	}

	var config *proxy.Config
	config, err = proxy.NewConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}

	reqFilter := proxy.New(config)
	err = (*http.Server)(reqFilter).ListenAndServe()
	log.Println(err)
}

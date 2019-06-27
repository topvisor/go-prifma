package proxy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfig_LoadFromFile(t *testing.T) {
	filename := "test.json"
	testJson := map[string]interface{} {
		"errorLogFile": "error.log",
		"listenPort": 8888,
		"listenType": "http",
		"filters": []interface{} {
			map[string]interface{} {
				"proxy": map[string]interface{} {
					"url": "http://example.com:3128",
					"connectHeaders": map[string]interface{} {
						"Proxy-Connection": "keep-alive",
					},
				},
				"enabled": true,
			},
			map[string]interface{} {
				"block": false,
				"enabled": true,
			},
		},
	}

	jsonData, err := json.Marshal(testJson)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(filename, jsonData, os.ModeTemporary)
	if err != nil {
		t.Fatal(err)
	}

	config := Config{}
	err = config.LoadFromFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(filename)
	if err != nil {
		t.Fatal(err)
	}
}
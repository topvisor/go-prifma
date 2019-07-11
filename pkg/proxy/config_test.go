package proxy

import (
	"encoding/json"
	"strconv"
	"testing"
)

var configMap = map[string]interface{}{
	"server": map[string]interface{}{
		"listenPort": 3128,
		"listenType": "http",
	},
	"accessLog":         "/path/to/access.log",
	"enableUseIpHeader": false,
}

func TestParseConfig(t *testing.T) {
	jsonData, err := json.Marshal(configMap)
	if err != nil {
		t.Fatal(err)
	}

	jsonStr := string(jsonData)
	config, err := ParseConfig(jsonStr)
	if err != nil {
		t.Fatal(err.Error())
	}

	listenPort := configMap["server"].(map[string]interface{})["listenPort"].(int)
	if config.Server.ListenPort != listenPort {
		t.Fatal("server.listenPort must be " + string(listenPort))
	}

	listenType := configMap["server"].(map[string]interface{})["listenType"].(string)
	if config.Server.ListenType != listenType {
		t.Fatal("server.listenPort must be " + listenType)
	}

	accessLog := configMap["accessLog"].(string)
	if *config.AccessLog != accessLog {
		t.Fatal("accessLog must be " + accessLog)
	}

	enableUseIpHeader := configMap["enableUseIpHeader"].(bool)
	if *config.EnableUseIpHeader != enableUseIpHeader {
		t.Fatalf("enableUseIpHeader must be " + strconv.FormatBool(enableUseIpHeader))
	}

	if config.BlockRequests != nil {
		t.Fatalf("blockRequests must be nil")
	}
}

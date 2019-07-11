package go_proxy_server

//import (
//	"encoding/json"
//	"io/ioutil"
//	"net/http"
//	"os"
//	"testing"
//)
//
//func TestNew(t *testing.T) {
//	filename := "test.json"
//	testJson := map[string]interface{}{
//		"listenPort": 9999,
//		"listenType": "http",
//		"filters": []interface{}{
//			map[string]interface{}{
//				"proxy": map[string]interface{}{
//					"url": "http://example.com:3128",
//					"connectHeaders": map[string]interface{}{
//						"Proxy-Connection": "keep-alive",
//					},
//				},
//				"enabled": true,
//			},
//		},
//	}
//
//	jsonData, err := json.Marshal(testJson)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = ioutil.WriteFile(filename, jsonData, os.ModeTemporary)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	var config *Config
//	config, err = NewConfig(filename)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = os.Remove(filename)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	reqfilter := New(config)
//	err = (*http.Server)(reqfilter).ListenAndServe()
//	t.Log(err)
//}

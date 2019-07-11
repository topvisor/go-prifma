package go_proxy_server

import (
	"encoding/json"
	"reflect"
)

type ListenType int

const (
	ListenTypeHttp ListenType = iota
)

func (t ListenType) MarshalJSON() ([]byte, error) {
	switch t {

	case ListenTypeHttp:
		return json.Marshal("http")

	default:
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(t),
			Str:   "invalid ListenType",
		}

	}
}

func (t *ListenType) UnmarshalJSON(data []byte) error {
	var stringListenType string

	err := json.Unmarshal(data, &stringListenType)
	if err != nil {
		return err
	}

	switch stringListenType {

	case "http":
		*t = ListenTypeHttp
		return nil

	default:
		return &json.UnmarshalTypeError{
			Value: "string",
			Type:  reflect.TypeOf(*t),
		}

	}
}

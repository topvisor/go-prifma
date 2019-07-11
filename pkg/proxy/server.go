package proxy

import (
	"errors"
	"net/http"
)

type listenType int

const (
	ListenTypeHttp listenType = iota
)

func ListenTypeFromString(lTypeStr string) (listenType, error) {
	switch lTypeStr {
	case "http":
		return ListenTypeHttp, nil
	default:
		return -1, errors.New("unavailable listen type")
	}
}

type Server struct {
	ListenPort int
	ListenType listenType
	Handler    Handler

	httpServer http.Server
}

func (t *Server) SetFromConfig(config Config) error {
	port := config.Server.ListenPort
	ltype, err := ListenTypeFromString(config.Server.ListenType)
	if err != nil {
		return err
	}

	handler := Handler{}
	if err = handler.setFromConfig(config.ConfigHandler); err != nil {
		return err
	}
	if err = t.Handler.Close(); err != nil {
		return err
	}

	t.ListenPort = port
	t.ListenType = ltype
	t.Handler = handler

	return nil
}

func (t *Server) LoadFromConfig(filename string) error {
	config, err := ParseConfigFromFile(filename)
	if err != nil {
		return err
	}

	if err = t.SetFromConfig(*config); err != nil {
		return err
	}

	return nil
}

func (t *Server) ListenAndServe() {
	err := t.httpServer.ListenAndServe()
	t.Handler.ErrorLogger.Fatalln(err)
}

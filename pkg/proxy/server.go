package proxy

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"net/http"
	"sync"
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

	httpServer            *http.Server
	fsWatcher             *fsnotify.Watcher
	watchedConfigFilename string
	loadFromConfigMutex   sync.Mutex
}

func (t *Server) SetFromConfig(config Config) error {
	port := config.Server.ListenPort
	ltype, err := ListenTypeFromString(config.Server.ListenType)
	if err != nil {
		return err
	}

	handler := Handler{}
	if err = handler.SetFromConfig(config.ConfigHandler); err != nil {
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
	t.loadFromConfigMutex.Lock()

	config, err := ParseConfigFromFile(filename)
	if err != nil {
		return err
	}
	if err = t.SetFromConfig(*config); err != nil {
		return err
	}

	t.loadFromConfigMutex.Unlock()

	return nil
}

func (t *Server) WatchForConfig(filename string) error {
	var err error

	if filename == t.watchedConfigFilename {
		return nil
	}
	if err = t.LoadFromConfig(filename); err != nil {
		return err
	}
	if t.fsWatcher != nil {
		if err = t.fsWatcher.Remove(t.watchedConfigFilename); err != nil {
			return err
		}
	} else {
		if t.fsWatcher, err = fsnotify.NewWatcher(); err != nil {
			return err
		}

		go t.listenFsWatcherErrors()
		go t.listenFsWatcherEvents()
	}

	if err = t.fsWatcher.Add(filename); err != nil {
		return err
	}

	t.watchedConfigFilename = filename

	return nil
}

func (t *Server) ListenAndServe() {
	err := t.httpServer.ListenAndServe()
	t.Handler.ErrorLogger.Fatalln(err)
}

func (t *Server) listenFsWatcherErrors() {
	for err := range t.fsWatcher.Errors {
		t.Handler.ErrorLogger.Println(err)
	}
}

func (t *Server) listenFsWatcherEvents() {
	for event := range t.fsWatcher.Events {
		if event.Name == t.watchedConfigFilename && event.Op == fsnotify.Write {
			if err := t.LoadFromConfig(event.Name); err != nil {
				t.Handler.ErrorLogger.Println(err)
			}
		}
	}
}

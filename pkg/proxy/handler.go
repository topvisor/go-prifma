package proxy

import (
	"log"
	"os"
)

type authType int

const (
	AuthTypeBasic authType = iota
)

type Handler struct {
	AccessLog            *log.Logger
	ErrorLog             *log.Logger
	AuthType             authType
	Htpasswd             *string
	HtpasswdForRedirects *string
}

func (t *Handler) SetFromConfig(config ConfigHandler) error {

	return nil
}

func (t *Handler) logErrorln(err error, isFatal bool) {
	if t.ErrorLog != nil {
		t.ErrorLog.Println(err)
	}
	if isFatal {
		os.Exit(1)
	}
}

func (t *Handler) logAccessln(access string) {
	if t.AccessLog != nil {
		t.AccessLog.Println(access)
	}
}

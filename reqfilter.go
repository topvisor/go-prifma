package proxyserver

import (
	"fmt"
	"net/http"
)

type ReqFilter http.Server

func New(config *Config) *ReqFilter {
	var reqfilter ReqFilter

	reqfilter.Addr = fmt.Sprintf(":%d", config.ListenPort)
	reqfilter.Handler = NewHandler(config.Filters)

	return &reqfilter
}

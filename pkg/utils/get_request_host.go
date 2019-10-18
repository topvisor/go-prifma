package utils

import (
	"net"
	"net/http"
)

func GetRequestHostname(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		host = req.Host
	}

	return host
}

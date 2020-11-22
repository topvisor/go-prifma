package utils

import (
	"net"
)

func GetHostname(host string) string {
	hostname, _, err := net.SplitHostPort(host)
	if err != nil {
		hostname = host
	}

	return hostname
}

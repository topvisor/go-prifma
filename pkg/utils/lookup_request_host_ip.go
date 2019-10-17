package utils

import (
	"net"
	"net/http"
)

func LookupIp(host string) (ipV4 net.IP, ipV6 net.IP, err error) {
	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		return nil, nil, err
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, nil, err
	}

	for _, ip := range ips {
		if ipV4Tmp := ip.To4(); ipV4Tmp != nil {
			if ipV4 == nil {
				ipV4 = ipV4Tmp
			}
		} else {
			if ipV6 == nil {
				ipV6 = ip
			}
		}

		if ipV4 != nil && ipV6 != nil {
			break
		}
	}

	return ipV4, ipV6, nil
}

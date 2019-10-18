package utils

import (
	"net/http"
	"regexp"
)

var ProxyHeadersRegexp, _ = regexp.Compile("^(?i)prifma-")

func RemoveProxyHeaders(req *http.Request) {
	for key := range req.Header {
		if ProxyHeadersRegexp.MatchString(key) {
			req.Header.Del(key)
		}
	}
}

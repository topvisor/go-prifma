package prifma

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const closeTimeout = time.Second * 2

func contextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	} else {
		return context.WithCancel(ctx)
	}
}

func closeFile(closer io.Closer) {
	if err := closer.Close(); err != nil && err != os.ErrClosed && err != os.ErrNotExist {
		time.Sleep(closeTimeout)
		_ = closer.Close()
	}
}

func proxyBasicAuth(req *http.Request) (username, password string, ok bool) {
	auth := req.Header.Get("Proxy-Authorization")
	if auth == "" {
		return
	}

	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

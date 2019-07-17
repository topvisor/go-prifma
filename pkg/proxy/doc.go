// Package proxy implements simple proxy server.
// Server supports a http proxy server type.
// Request's handler is selected based on the given conditions. It can authenticate users using basic http
// authentication, select outgoing ip, send requests through another proxy server and block requests.
package proxy

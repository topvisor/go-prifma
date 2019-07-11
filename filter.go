package go_proxy_server

type Filter struct {
	Proxy   *Proxy `json:"proxy,omitempty"`
	Block   bool   `json:"block,omitempty"`
	Enabled bool   `json:"enabled"`
}

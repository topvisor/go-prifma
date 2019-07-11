package proxy

type Proxy struct {
	htpasswdForRedirects *BasicAuth
}

func (t *Proxy) setFromConfig(config ConfigProxy) error {
	return nil
}

package proxy

import "net/url"

type Proxy struct {
	Url          url.URL
	Htpasswd     *BasicAuth
	ProxyHeaders map[string]string
}

func (t *Proxy) setFromConfig(config ConfigProxy) error {
	parserUrl, err := url.ParseRequestURI(config.Url)
	if err != nil {
		return err
	}

	var htpasswd *BasicAuth
	if config.Htpasswd != nil {
		htpasswd, err = NewBasicAuth(*config.Htpasswd)
		if err != nil {
			return err
		}
	}

	t.Url = *parserUrl
	t.Htpasswd = htpasswd
	t.ProxyHeaders = config.ProxyHeaders

	return nil
}

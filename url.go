package go_proxy_server

import (
	"encoding/json"
	"net/url"
)

type URL url.URL

func (t *URL) MarshalJSON() ([]byte, error) {
	return json.Marshal((*url.URL)(t).String())
}

func (t *URL) UnmarshalJSON(data []byte) error {
	var rawurl string
	err := json.Unmarshal(data, &rawurl)
	if err != nil {
		return err
	}

	var parsedUrl *url.URL
	parsedUrl, err = url.Parse(rawurl)
	if err != nil {
		return err
	}

	*t = URL(*parsedUrl)

	return nil
}

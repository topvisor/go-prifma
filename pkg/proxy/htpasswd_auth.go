package proxy

import (
	"encoding/csv"
	auth "github.com/abbot/go-http-auth"
	"os"
)

// Htpasswd is an authenticator implementation for basic http authentication.
// It uses Htpasswd of "github.com/abbot/go-http-auth"
type Htpasswd struct {
	users map[string]string
	auth.BasicAuth
}

// LoadHtpasswd load ".BasicAuth" file
func LoadHtpasswd(filename string) (*Htpasswd, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	reader.Comma = ':'
	reader.Comment = '#'
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	htpasswdAuth := new(Htpasswd)
	htpasswdAuth.Secrets = htpasswdAuth.secrets
	htpasswdAuth.Headers = auth.ProxyHeaders
	htpasswdAuth.users = make(map[string]string)
	for _, record := range records {
		htpasswdAuth.users[record[0]] = record[1]
	}

	return htpasswdAuth, nil
}

// secrets is implementation of SecretProvider of "github.com/abbot/go-http-auth"
func (t *Htpasswd) secrets(user, realm string) string {
	password, exists := t.users[user]
	if !exists {
		return ""
	}

	return password
}

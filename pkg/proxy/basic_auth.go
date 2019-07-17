package proxy

import (
	"encoding/csv"
	auth "github.com/abbot/go-http-auth"
	"net/http"
	"os"
)

// BasicAuth is an authenticator implementation for basic http authentication.
// It uses BasicAuth of "github.com/abbot/go-http-auth"
type BasicAuth struct {
	Users map[string]string

	basicAuth auth.BasicAuth
}

// NewBasicAuth load ".htpasswd" file
func NewBasicAuth(filename string) (*BasicAuth, error) {
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

	basicAuth := new(BasicAuth)
	basicAuth.Users = make(map[string]string)
	for _, record := range records {
		basicAuth.Users[record[0]] = record[1]
	}

	return basicAuth, nil
}

// CheckAuth checks the username/password combination from the
// request. Returns either an empty string (authentication failed) or
// the name of the authenticated user.
func (t *BasicAuth) CheckAuth(r *http.Request) string {
	t.initIfNeed()

	return t.basicAuth.CheckAuth(r)
}

// RequireAuth is an http.HandlerFunc for BasicAuth which initiates
// the authentication process (or requires reauthentication).
func (t *BasicAuth) RequireAuth(w http.ResponseWriter, r *http.Request) {
	t.initIfNeed()

	t.basicAuth.RequireAuth(w, r)
}

// initIfNeed initiates BasicAuth for use if it's not initiated
func (t *BasicAuth) initIfNeed() {
	if t.basicAuth.Secrets == nil {
		t.basicAuth.Secrets = t.secrets
		t.basicAuth.Headers = auth.ProxyHeaders
	}
}

// secrets is implementation of SecretProvider of "github.com/abbot/go-http-auth"
func (t *BasicAuth) secrets(user, realm string) string {
	password, exists := t.Users[user]
	if !exists {
		return ""
	}

	return password
}

package proxy

import (
	"encoding/csv"
	auth "github.com/abbot/go-http-auth"
	"golang.org/x/net/context"
	"net/http"
	"os"
)

type BasicAuth struct {
	Users map[string]string

	basicAuth auth.BasicAuth
}

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

func (t *BasicAuth) CheckAuth(r *http.Request) string {
	t.initIfNeed()

	return t.basicAuth.CheckAuth(r.Request)
}

func (t *BasicAuth) NewContext(ctx context.Context, r *http.Request) context.Context {
	t.initIfNeed()

	return t.basicAuth.NewContext(ctx, r.Request)
}

func (t *BasicAuth) RequireAuth(w http.ResponseWriter, r *http.Request) {
	t.initIfNeed()

	t.basicAuth.RequireAuth(w, r.Request)
}

func (t *BasicAuth) initIfNeed() {
	if t.basicAuth.Secrets == nil {
		t.basicAuth.Secrets = t.sercers
		t.basicAuth.Headers = auth.ProxyHeaders
	}
}

func (t *BasicAuth) sercers(user, realm string) string {
	password, exists := t.Users[user]
	if !exists {
		return ""
	}

	return password
}

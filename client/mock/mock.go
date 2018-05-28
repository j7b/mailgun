// Package mock implements a mock client for testing.
package mock

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/j7b/mailgun/client"
)

type rt struct {
	*testing.T
	http.RoundTripper
}

func (rt *rt) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.Log("method", req.Method, "url", req.URL.String())
	req.URL.Path = path.Join(req.URL.Path, req.Method)
	for k, v := range req.Header {
		for _, v := range v {
			rt.Log("header", k, v)
		}
	}
	return rt.RoundTripper.RoundTrip(req)
}

func mock(t *testing.T) *http.Client {
	gp := os.Getenv(`GOPATH`)
	if len(gp) == 0 {
		t.Fatal(`GOPATH not set`)
	}
	if l := len(filepath.SplitList(gp)); l != 1 {
		t.Fatal(`GOPATH appears split`)
	}
	mockdir := filepath.Join(gp, `src`, `github.com`, `j7b`, `mailgun`, `client`, `mock`)
	fi, err := os.Stat(mockdir)
	if err != nil {
		t.Fatal(err)
	}
	if !fi.IsDir() {
		t.Fatal(mockdir + ` not directory`)
	}
	dir := http.Dir(mockdir)
	c := new(http.Client)
	c.Transport = &rt{RoundTripper: http.NewFileTransport(dir), T: t}
	return c
}

func Client(t *testing.T) client.Caller {
	c := &client.Requester{}
	c.Endpoint = `_mock/`
	c.APIDomain = `domain.mock`
	h := mock(t)
	c.Client = h
	return c
}

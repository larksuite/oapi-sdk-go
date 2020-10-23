package transport

import (
	"net/http"
)

var DefaultClient *http.Client

func init() {
	transport := &Transport{}
	DefaultClient = transport.Client()
}

type Transport struct {
}

func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.transport().RoundTrip(req)
}

func (t *Transport) transport() http.RoundTripper {
	return http.DefaultTransport
}

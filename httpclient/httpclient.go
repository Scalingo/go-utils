package httpclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
)

type ClientOpt func(c *http.Client)

func WithTimeout(d time.Duration) ClientOpt {
	return func(c *http.Client) {
		c.Timeout = d
	}
}

func WithTLSConfig(config *tls.Config) ClientOpt {
	return func(c *http.Client) {
		parent := &http.Transport{
			TLSClientConfig: config,
		}
		c.Transport = reqidTransport{parent: parent}
	}
}

func WithAuthentication(username, password string) ClientOpt {
	return func(c *http.Client) {
		c.Transport = authTransport{
			parent:   reqidTransport{http.DefaultTransport},
			username: username, password: password,
		}
	}
}

func NewClient(opts ...ClientOpt) *http.Client {
	client := &http.Client{
		Transport: reqidTransport{parent: http.DefaultTransport},
	}
	for _, o := range opts {
		o(client)
	}
	return client
}

type reqidTransport struct {
	parent http.RoundTripper
}

func (t reqidTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("X-Request-ID") != "" {
		return t.parent.RoundTrip(req)
	}

	reqID, ok := req.Context().Value("request_id").(string)
	if !ok {
		uuid, err := uuid.NewV4()
		if err != nil {
			return nil, fmt.Errorf("fail to generate UUID for X-Request-ID: %v", err)
		}
		reqID = uuid.String()
	}
	req.Header.Set("X-Request-ID", reqID)
	return t.parent.RoundTrip(req)
}

type authTransport struct {
	parent   http.RoundTripper
	username string
	password string
}

func (t authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if _, _, ok := req.BasicAuth(); !ok {
		req.SetBasicAuth(t.username, t.password)
	}
	return t.parent.RoundTrip(req)
}

package httpclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

type ClientOpt func(c *client)

type client struct {
	config   *tls.Config
	user     string
	password string
	timeout  time.Duration
}

func WithTimeout(d time.Duration) ClientOpt {
	return func(c *client) {
		c.timeout = d
	}
}

func WithTLSConfig(config *tls.Config) ClientOpt {
	return func(c *client) {
		c.config = config
	}
}

func WithAuthentication(username, password string) ClientOpt {
	return func(c *client) {
		c.user = username
		c.password = password
	}
}

func NewClient(opts ...ClientOpt) *http.Client {
	httpClient := &http.Client{
		Transport: reqidTransport{parent: http.DefaultTransport},
	}
	c := client{}
	for _, o := range opts {
		o(&c)
	}
	if c.timeout > 0 {
		httpClient.Timeout = c.timeout
	}
	if c.config != nil {
		parent := &http.Transport{
			TLSClientConfig: c.config,
		}
		httpClient.Transport = reqidTransport{parent: parent}
	}
	if c.user != "" || c.password != "" {
		httpClient.Transport = authTransport{
			parent:   httpClient.Transport,
			username: c.user, password: c.password,
		}
	}
	return httpClient
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

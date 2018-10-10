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
		c.Transport = transport{parentTransport: parent}
	}
}

func NewClient(opts ...ClientOpt) *http.Client {
	client := &http.Client{
		Transport: transport{parentTransport: http.DefaultTransport},
	}
	for _, o := range opts {
		o(client)
	}
	return client
}

type transport struct {
	parentTransport http.RoundTripper
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("X-Request-ID") != "" {
		return t.parentTransport.RoundTrip(req)
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
	return t.parentTransport.RoundTrip(req)
}

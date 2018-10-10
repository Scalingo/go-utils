package httpclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cases := []struct {
		Name    string
		Context context.Context
		Expect  func(*testing.T, string)
	}{
		{
			Name:    "it should add a X-Request-ID header if present in request context",
			Context: context.WithValue(context.Background(), "request_id", "123"),
			Expect: func(t *testing.T, body string) {
				assert.Equal(t, "123", body)
			},
		}, {
			Name:    "a UUID should be added of no request ID is in the context",
			Context: context.Background(),
			Expect: func(t *testing.T, body string) {
				assert.Len(t, body, 36)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, r.Header.Get("X-Request-ID"))
			}))
			defer server.Close()

			req, err := http.NewRequest("GET", server.URL, nil)
			require.NoError(t, err)
			req = req.WithContext(c.Context)
			client := NewClient()
			res, err := client.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			c.Expect(t, string(body))
		})
	}
}

func TestNewClient_WithAuthentication(t *testing.T) {
	cases := []struct {
		Name    string
		Request func(*testing.T, *http.Request) *http.Request
		Expect  func(*testing.T, string)
	}{
		{
			Name: "it should add basic authentication to the request",
			Request: func(t *testing.T, req *http.Request) *http.Request {
				return req
			},
			Expect: func(t *testing.T, body string) {
				assert.Equal(t, "user:password", body)
			},
		}, {
			Name: "it should not add the authentication if already set",
			Request: func(t *testing.T, req *http.Request) *http.Request {
				req.SetBasicAuth("pipo", "molo")
				return req
			},
			Expect: func(t *testing.T, body string) {
				assert.Equal(t, "pipo:molo", body)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user, pass, _ := r.BasicAuth()
				fmt.Fprintf(w, "%v:%v", user, pass)
			}))
			defer server.Close()

			req, err := http.NewRequest("GET", server.URL, nil)
			require.NoError(t, err)
			req = c.Request(t, req)
			client := NewClient(WithAuthentication("user", "password"))
			res, err := client.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			c.Expect(t, string(body))
		})
	}
}

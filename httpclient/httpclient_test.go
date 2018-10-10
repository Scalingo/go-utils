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
		Name   string
		Expect func(t *testing.T, url string)
	}{
		{
			Name: "it should add a X-Request-ID header if present in request context",
			Expect: func(t *testing.T, url string) {
				ctx := context.WithValue(context.Background(), "request_id", "123")
				req, err := http.NewRequest("GET", url, nil)
				require.NoError(t, err)
				req = req.WithContext(ctx)
				c := NewClient()
				res, err := c.Do(req)
				require.NoError(t, err)
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, "123", string(body))
			},
		}, {
			Name: "a UUID should be added of no request ID is in the context",
			Expect: func(t *testing.T, url string) {
				req, err := http.NewRequest("GET", url, nil)
				require.NoError(t, err)
				c := NewClient()
				res, err := c.Do(req)
				require.NoError(t, err)
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Len(t, string(body), 36)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, r.Header.Get("X-Request-ID"))
			}))
			defer server.Close()
			c.Expect(t, server.URL)
		})
	}
}

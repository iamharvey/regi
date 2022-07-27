package rest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultServerURL(t *testing.T) {
	cases := []struct {
		name       string
		scheme     string
		host       string
		apiPrefix  string
		defaultTLS bool
		err        error
	}{
		{
			name:       "valid URL",
			scheme:     "http",
			host:       "api.server",
			apiPrefix:  "/v1",
			defaultTLS: false,
			err:        nil,
		},
		{
			name:       "valid URL with port",
			scheme:     "https",
			host:       "api.server:8080",
			apiPrefix:  "/v1",
			defaultTLS: true,
			err:        nil,
		},
		{
			name:       "scheme and TLS setting does not match",
			scheme:     "http",
			host:       "api.server/path/to/visit",
			apiPrefix:  "/v1",
			defaultTLS: false,
			err:        fmt.Errorf("host must be a URL or a host:port pair: %q", "api.server/path/to/visit"),
		},
	}
	for _, c := range cases {
		url, _, err := DefaultServerURL(c.host, c.apiPrefix, c.defaultTLS)
		assert.Equal(t, err, c.err)
		if url != nil {
			assert.Equal(t, c.scheme, url.Scheme)
			assert.Equal(t, c.host, url.Host)
		}
	}
}

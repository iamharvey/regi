package rest

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"time"
)

// Client defines a common JAC CMS API client.
type Client struct {
	// base is the root URL for all invocations of the client.
	base *url.URL

	// versionedAPIPath is a path segment connecting the base URL to the resource root
	versionedAPIPath string

	contentConfig *ContentConfig

	*http.Client
}

// NewClient returns a API server client, which is an HTTP client.
func NewClient(cfg *ClientConfig) (*Client, error) {
	var (
		client *http.Client
		err    error
	)

	enableTLS := false
	// Make an HTTPS client if TLS settings are available.
	tlsConfig := cfg.TLSClientConfig
	if tlsConfig != nil {
		client, err = createHTTPClientWithTLS(tlsConfig)
		if err != nil {
			return nil, err
		}
		enableTLS = true
	}

	base, versionedAPIPath, err := DefaultServerURL(cfg.Host, cfg.APIPath, enableTLS)
	if err != nil {
		return nil, err
	}

	// SetCurrentContext client timeout.
	client = &http.Client{}
	client.Timeout = time.Second * 2

	return &Client{
		base:             base,
		versionedAPIPath: versionedAPIPath,
		Client:           client,
		contentConfig:    cfg.ContentConfig,
	}, nil
}

// Verb receives a verb which indicates what HTTP method should be invoked.
func (c *Client) Verb(verb string) *Request {
	return NewRequest(c).Verb(verb)
}

// createHTTPClientWithTLS creates an HTTPS client.
func createHTTPClientWithTLS(tlsConfig *TLSClientConfig) (*http.Client, error) {
	cert, err := tls.X509KeyPair(tlsConfig.CertData, tlsConfig.KeyData)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(tlsConfig.CAData)

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			},
		},
	}, nil
}

package rest

import (
	"fmt"
	"net/url"
	"path"
)

// DefaultServerURL converts a host, host:port, or URL string to the default base server API path
// to use with a Client at a given API version following the standard conventions for a
// JAC CMS API.
func DefaultServerURL(host, apiPrefix string, defaultTLS bool) (*url.URL, string, error) {
	if host == "" {
		return nil, "", fmt.Errorf("host must be a URL or a host:port pair")
	}
	base := host
	hostURL, err := url.Parse(base)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		scheme := "http://"
		if defaultTLS {
			scheme = "https://"
		}

		hostURL, err = url.Parse(scheme + base)
		if err != nil {
			return nil, "", err
		}

		if hostURL.Path != "" && hostURL.Path != "/" {
			return nil, "", fmt.Errorf("host must be a URL or a host:port pair: %q", base)
		}
	}

	versionedAPIPath := path.Join("/", apiPrefix)

	return hostURL, versionedAPIPath, nil
}

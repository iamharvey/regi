package rest

import (
	"fmt"
	"time"
)

// ClientConfig defines a config for creating a REST client.
type ClientConfig struct {
	// Host must be a host string, a host:port pair, or a URL to the base of the apiserver.
	Host string

	// APIPath is a sub-path that points to an API root.
	APIPath string

	// ContentConfig contains settings that affect how objects are transformed when
	// sent to the server.
	ContentConfig *ContentConfig

	// BearerToken is needed when server requires Bearer authentication (not refreshable).
	BearerToken string

	// BearerTokenTag defines what tag name to be used.
	// If not set, default 'bearer' will be set.
	BearerTokenTag string

	TLSClientConfig *TLSClientConfig

	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout time.Duration
}

// ContentConfig contains settings that affect how objects are transformed when
// sent to the server.
type ContentConfig struct {
	// AcceptContentTypes specifies the types the client will accept and is optional.
	// If not set, ContentType will be used to define the Accept header
	AcceptContentTypes string

	// ContentType will be set as the Accept header on requests made to the server, and
	// as the default content type on any object sent to the server. If not set,
	// "application/json" is used.
	ContentType string
}

// TLSClientConfig contains settings to enable transport layer security (TLS).
type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	Insecure bool

	// ServerName is passed to the server for SNI and is used in the client to check server
	// certificates against. If ServerName is empty, the hostname used to contact the
	// server is used.
	ServerName string

	// Server requires TLS client certificate authentication
	CertFile string

	// Server requires TLS client certificate authentication
	KeyFile string

	// Trusted root certificates for server
	CAFile string

	// CertData holds PEM-encoded bytes (typically read from a client certificate file).
	// CertData takes precedence over CertFile
	CertData []byte

	// KeyData holds PEM-encoded bytes (typically read from a client certificate key file).
	// KeyData takes precedence over KeyFile
	KeyData []byte

	// CAData holds PEM-encoded bytes (typically read from a root certificates bundle).
	// CAData takes precedence over CAFile
	CAData []byte
}

var _ fmt.Stringer = TLSClientConfig{}
var _ fmt.GoStringer = TLSClientConfig{}

type sanitizedTLSClientConfig TLSClientConfig

// GoString implements fmt.GoStringer and sanitizes sensitive fields of
// TLSClientConfig to prevent accidental leaking via logs.
func (c TLSClientConfig) GoString() string {
	return c.String()
}

// String implements fmt.Stringer and sanitizes sensitive fields of
// TLSClientConfig to prevent accidental leaking via logs.
func (c TLSClientConfig) String() string {
	cc := sanitizedTLSClientConfig{
		Insecure:   c.Insecure,
		ServerName: c.ServerName,
		CertFile:   c.CertFile,
		KeyFile:    c.KeyFile,
		CAFile:     c.CAFile,
		CertData:   c.CertData,
		KeyData:    c.KeyData,
		CAData:     c.CAData,
	}

	// Explicitly mark non-empty credential fields as redacted.
	if len(cc.CertData) != 0 {
		cc.CertData = []byte("--- TRUNCATED ---")
	}
	if len(cc.KeyData) != 0 {
		cc.KeyData = []byte("--- REDACTED ---")
	}
	return fmt.Sprintf("%#v", cc)
}
